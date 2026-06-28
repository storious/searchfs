package client

import (
	"context"
	"io"
	"strings"
	"testing"

	"gdfs/internal/datanode"

	"github.com/stretchr/testify/require"
)

type fakeBlockReader struct {
	blocks map[datanode.BlockID]string
}

func newFakeBlockReader(blocks map[datanode.BlockID]string) *fakeBlockReader {
	return &fakeBlockReader{blocks: blocks}
}

func (r *fakeBlockReader) GetBlock(ctx context.Context, id datanode.BlockID) (io.ReadCloser, error) {
	data, ok := r.blocks[id]
	if !ok {
		return nil, io.ErrUnexpectedEOF
	}

	return io.NopCloser(strings.NewReader(data)), nil
}

func TestReaderReconstructsBlocks(t *testing.T) {
	blocks := []datanode.BlockInfo{
		{ID: "block-001", Size: 5},
		{ID: "block-002", Size: 5},
		{ID: "block-003", Size: 1},
	}

	source := newFakeBlockReader(map[datanode.BlockID]string{
		"block-001": "hello",
		"block-002": "-worl",
		"block-003": "d",
	})

	reader, err := NewReader(source)
	require.NoError(t, err)

	var out strings.Builder

	n, err := reader.Read(context.Background(), blocks, &out)
	require.NoError(t, err)

	require.Equal(t, int64(len("hello-world")), n)
	require.Equal(t, "hello-world", out.String())
}

func TestNewReaderRejectsNilSource(t *testing.T) {
	reader, err := NewReader(nil)

	require.Error(t, err)
	require.Nil(t, reader)
}
