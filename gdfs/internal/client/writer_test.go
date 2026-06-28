package client

import (
	"context"
	"io"
	"strings"
	"testing"

	"gdfs/internal/datanode"

	"github.com/stretchr/testify/require"
)

type fakeBlockWriter struct {
	blocks map[datanode.BlockID]string
}

func newFakeBlockWriter() *fakeBlockWriter {
	return &fakeBlockWriter{
		blocks: make(map[datanode.BlockID]string),
	}
}

func (w *fakeBlockWriter) PutBlock(ctx context.Context, id datanode.BlockID, r io.Reader) (datanode.BlockInfo, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return datanode.BlockInfo{}, err
	}

	w.blocks[id] = string(data)

	return datanode.BlockInfo{
		ID:       id,
		Size:     int64(len(data)),
		Checksum: "fake-checksum",
	}, nil
}

func TestWriterSplitsInputIntoBlocks(t *testing.T) {
	target := newFakeBlockWriter()

	writer, err := NewWriter(5, target)
	require.NoError(t, err)

	result, err := writer.Write(context.Background(), strings.NewReader("hello-world"))
	require.NoError(t, err)

	require.Equal(t, int64(len("hello-world")), result.Size)
	require.Len(t, result.Blocks, 3)

	require.Contains(t, values(target.blocks), "hello")
	require.Contains(t, values(target.blocks), "-worl")
	require.Contains(t, values(target.blocks), "d")
}

func values(m map[datanode.BlockID]string) []string {
	out := make([]string, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	return out
}
