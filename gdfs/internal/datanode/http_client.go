package datanode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HTTPClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		baseURL: baseURL,
		client:  http.DefaultClient,
	}
}

func (c *HTTPClient) PutBlock(ctx context.Context, id BlockID, r io.Reader) (BlockInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.blockURL(id), r)
	if err != nil {
		return BlockInfo{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return BlockInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return BlockInfo{}, fmt.Errorf("put block failed: status=%s", resp.Status)
	}

	var info BlockInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return BlockInfo{}, err
	}

	return info, nil
}

func (c *HTTPClient) GetBlock(ctx context.Context, id BlockID) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.blockURL(id), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("get block failed: status=%s", resp.Status)
	}

	return resp.Body, nil
}

func (c *HTTPClient) DeleteBlock(ctx context.Context, id BlockID) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.blockURL(id), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete block failed: status=%s", resp.Status)
	}

	return nil
}

func (c *HTTPClient) StatBlock(ctx context.Context, id BlockID) (BlockInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, c.blockURL(id), nil)
	if err != nil {
		return BlockInfo{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return BlockInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return BlockInfo{}, fmt.Errorf("stat block failed: status=%s", resp.Status)
	}

	return BlockInfo{
		ID:       id,
		Size:     parseInt64(resp.Header.Get("X-Block-Size")),
		Checksum: resp.Header.Get("X-Block-Checksum"),
	}, nil
}

func (c *HTTPClient) blockURL(id BlockID) string {
	return c.baseURL + "/blocks/" + string(id)
}

func parseInt64(s string) int64 {
	var n int64
	_, _ = fmt.Sscanf(s, "%d", &n)
	return n
}

// helper for tests or future in-memory use
func readerFromBytes(b []byte) io.Reader {
	return bytes.NewReader(b)
}
