package client

import (
	"context"
	"io"

	"gdfs/internal/namenode"
)

type MetadataClient interface {
	PutFile(ctx context.Context, meta namenode.FileMetadata) (namenode.FileMetadata, error)
	GetFile(ctx context.Context, path namenode.FilePath) (namenode.FileMetadata, error)
	DeleteFile(ctx context.Context, path namenode.FilePath) error
}

type FileClient struct {
	writer   *Writer
	reader   *Reader
	metadata MetadataClient
}

type BlockClient interface {
	BlockWriter
	BlockReader
}

func NewFileClient(blockSize int64, blocks BlockClient, metadata MetadataClient) (*FileClient, error) {
	if blocks == nil {
		return nil, ErrNilBlockClient
	}
	if metadata == nil {
		return nil, ErrNilMetadataClient
	}

	writer, err := NewWriter(blockSize, blocks)
	if err != nil {
		return nil, err
	}

	reader, err := NewReader(blocks)
	if err != nil {
		return nil, err
	}

	return &FileClient{
		writer:   writer,
		reader:   reader,
		metadata: metadata,
	}, nil
}

func (c *FileClient) PutFile(ctx context.Context, path namenode.FilePath, r io.Reader) (namenode.FileMetadata, error) {
	result, err := c.writer.Write(ctx, r)
	if err != nil {
		return namenode.FileMetadata{}, err
	}

	meta := namenode.FileMetadata{
		Path:   path,
		Size:   result.Size,
		Blocks: result.Blocks,
	}

	return c.metadata.PutFile(ctx, meta)
}

func (c *FileClient) GetFile(ctx context.Context, path namenode.FilePath, dst io.Writer) (int64, error) {
	meta, err := c.metadata.GetFile(ctx, path)
	if err != nil {
		return 0, err
	}

	return c.reader.Read(ctx, meta.Blocks, dst)
}

func (c *FileClient) StatFile(ctx context.Context, path namenode.FilePath) (namenode.FileMetadata, error) {
	return c.metadata.GetFile(ctx, path)
}

func (c *FileClient) DeleteFile(ctx context.Context, path namenode.FilePath) error {
	return c.metadata.DeleteFile(ctx, path)
}
