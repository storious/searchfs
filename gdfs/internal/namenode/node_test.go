package namenode

import (
	"context"
	"testing"

	"gdfs/internal/datanode"

	"github.com/stretchr/testify/require"
)

func TestNameNodeCreateGetDeleteFile(t *testing.T) {
	node, err := NewNameNode(NewMetadataStore())
	require.NoError(t, err)

	ctx := context.Background()

	meta := FileMetadata{
		Path: "/docs/hello.txt",
		Size: 11,
		Blocks: []datanode.BlockInfo{
			{ID: "block-001", Size: 5, Checksum: "a"},
			{ID: "block-002", Size: 6, Checksum: "b"},
		},
	}

	err = node.CreateFile(ctx, meta)
	require.NoError(t, err)

	require.True(t, node.ExistsFile(ctx, "/docs/hello.txt"))

	got, err := node.GetFile(ctx, "/docs/hello.txt")
	require.NoError(t, err)
	require.Equal(t, meta, got)

	err = node.DeleteFile(ctx, "/docs/hello.txt")
	require.NoError(t, err)

	require.False(t, node.ExistsFile(ctx, "/docs/hello.txt"))
}

func TestNewNameNodeRejectsNilStore(t *testing.T) {
	node, err := NewNameNode(nil)

	require.Error(t, err)
	require.Nil(t, node)
}
