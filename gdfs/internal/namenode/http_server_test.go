package namenode

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"gdfs/internal/datanode"

	"github.com/stretchr/testify/require"
)

func newTestHTTPServer(t *testing.T) *httptest.Server {
	t.Helper()

	node, err := NewNameNode(NewMetadataStore())
	require.NoError(t, err)

	return httptest.NewServer(NewHTTPServer(node))
}

func TestHTTPServerPutGetFile(t *testing.T) {
	server := newTestHTTPServer(t)
	defer server.Close()

	meta := FileMetadata{
		Size: 11,
		Blocks: []datanode.BlockInfo{
			{ID: "block-001", Size: 5, Checksum: "a"},
			{ID: "block-002", Size: 6, Checksum: "b"},
		},
	}

	body, err := json.Marshal(meta)
	require.NoError(t, err)

	putResp, err := http.DefaultClient.Do(mustRequest(
		t,
		http.MethodPut,
		server.URL+"/files/docs/hello.txt",
		bytes.NewReader(body),
	))
	require.NoError(t, err)
	defer putResp.Body.Close()

	require.Equal(t, http.StatusCreated, putResp.StatusCode)

	getResp, err := http.Get(server.URL + "/files/docs/hello.txt")
	require.NoError(t, err)
	defer getResp.Body.Close()

	require.Equal(t, http.StatusOK, getResp.StatusCode)

	var got FileMetadata
	err = json.NewDecoder(getResp.Body).Decode(&got)
	require.NoError(t, err)

	require.Equal(t, FilePath("/docs/hello.txt"), got.Path)
	require.Equal(t, meta.Size, got.Size)
	require.Equal(t, meta.Blocks, got.Blocks)
}

func TestHTTPServerDeleteFile(t *testing.T) {
	server := newTestHTTPServer(t)
	defer server.Close()

	body, err := json.Marshal(FileMetadata{
		Size: 5,
		Blocks: []datanode.BlockInfo{
			{ID: "block-001", Size: 5},
		},
	})
	require.NoError(t, err)

	putResp, err := http.DefaultClient.Do(mustRequest(
		t,
		http.MethodPut,
		server.URL+"/files/docs/hello.txt",
		bytes.NewReader(body),
	))
	require.NoError(t, err)
	defer putResp.Body.Close()

	require.Equal(t, http.StatusCreated, putResp.StatusCode)

	deleteResp, err := http.DefaultClient.Do(mustRequest(
		t,
		http.MethodDelete,
		server.URL+"/files/docs/hello.txt",
		nil,
	))
	require.NoError(t, err)
	defer deleteResp.Body.Close()

	require.Equal(t, http.StatusNoContent, deleteResp.StatusCode)

	getResp, err := http.Get(server.URL + "/files/docs/hello.txt")
	require.NoError(t, err)
	defer getResp.Body.Close()

	require.Equal(t, http.StatusNotFound, getResp.StatusCode)
}

func TestHTTPServerMissingFilePath(t *testing.T) {
	server := newTestHTTPServer(t)
	defer server.Close()

	resp, err := http.Get(server.URL + "/files/")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestHTTPServerMethodNotAllowed(t *testing.T) {
	server := newTestHTTPServer(t)
	defer server.Close()

	resp, err := http.DefaultClient.Do(mustRequest(
		t,
		http.MethodPost,
		server.URL+"/files/docs/hello.txt",
		nil,
	))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func mustRequest(t *testing.T, method, url string, body io.Reader) *http.Request {
	t.Helper()

	req, err := http.NewRequest(method, url, body)
	require.NoError(t, err)
	return req
}
