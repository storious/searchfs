package namenode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type HTTPClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  http.DefaultClient,
	}
}

func (c *HTTPClient) PutFile(ctx context.Context, meta FileMetadata) (FileMetadata, error) {
	body, err := json.Marshal(meta)
	if err != nil {
		return FileMetadata{}, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		c.fileURL(meta.Path),
		bytes.NewReader(body),
	)
	if err != nil {
		return FileMetadata{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return FileMetadata{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return FileMetadata{}, fmt.Errorf("put file failed: status=%s", resp.Status)
	}

	var out FileMetadata
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return FileMetadata{}, err
	}

	return out, nil
}

func (c *HTTPClient) GetFile(ctx context.Context, path FilePath) (FileMetadata, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.fileURL(path),
		nil,
	)
	if err != nil {
		return FileMetadata{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return FileMetadata{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return FileMetadata{}, fmt.Errorf("get file failed: status=%s", resp.Status)
	}

	var meta FileMetadata
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return FileMetadata{}, err
	}

	return meta, nil
}

func (c *HTTPClient) DeleteFile(ctx context.Context, path FilePath) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		c.fileURL(path),
		nil,
	)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete file failed: status=%s", resp.Status)
	}

	return nil
}

func (c *HTTPClient) fileURL(path FilePath) string {
	clean := strings.TrimPrefix(string(path), "/")
	return c.baseURL + "/files/" + clean
}
