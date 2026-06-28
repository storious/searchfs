package datanode

import (
	"context"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPClientPutGetStatDeleteBlock(t *testing.T) {
	store := NewLocalBlockStore(t.TempDir())

	node, err := NewDataNode("node-1", "127.0.0.1:0", store)
	require.NoError(t, err)

	server := httptest.NewServer(NewHTTPServer(node))
	defer server.Close()

	client := NewHTTPClient(server.URL)
	ctx := context.Background()

	info, err := client.PutBlock(ctx, BlockID("block-001"), strings.NewReader("hello client"))
	require.NoError(t, err)
	require.Equal(t, BlockID("block-001"), info.ID)
	require.Equal(t, int64(len("hello client")), info.Size)
	require.NotEmpty(t, info.Checksum)

	stat, err := client.StatBlock(ctx, BlockID("block-001"))
	require.NoError(t, err)
	require.Equal(t, info.ID, stat.ID)
	require.Equal(t, info.Size, stat.Size)
	require.Equal(t, info.Checksum, stat.Checksum)

	r, err := client.GetBlock(ctx, BlockID("block-001"))
	require.NoError(t, err)
	defer r.Close()

	body, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, "hello client", string(body))

	err = client.DeleteBlock(ctx, BlockID("block-001"))
	require.NoError(t, err)

	_, err = client.StatBlock(ctx, BlockID("block-001"))
	require.Error(t, err)
}
