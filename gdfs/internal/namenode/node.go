package namenode

import (
	"context"
	"errors"
)

type NameNode struct {
	store *MetadataStore
}

func NewNameNode(store *MetadataStore) (*NameNode, error) {
	if store == nil {
		return nil, errors.New("nil metadata store")
	}

	return &NameNode{store: store}, nil
}

func (n *NameNode) CreateFile(ctx context.Context, meta FileMetadata) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return n.store.PutFile(meta)
}

func (n *NameNode) GetFile(ctx context.Context, path FilePath) (FileMetadata, error) {
	select {
	case <-ctx.Done():
		return FileMetadata{}, ctx.Err()
	default:
	}

	return n.store.GetFile(path)
}

func (n *NameNode) DeleteFile(ctx context.Context, path FilePath) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return n.store.DeleteFile(path)
}

func (n *NameNode) ExistsFile(ctx context.Context, path FilePath) bool {
	select {
	case <-ctx.Done():
		return false
	default:
	}

	return n.store.Exists(path)
}
