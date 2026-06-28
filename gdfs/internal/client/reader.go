package client

import (
	"context"
	"io"

	"gdfs/internal/datanode"
)

type BlockReader interface {
	GetBlock(ctx context.Context, id datanode.BlockID) (io.ReadCloser, error)
}

type Reader struct {
	source BlockReader
}

func NewReader(source BlockReader) (*Reader, error) {
	if source == nil {
		return nil, ErrNilBlockReader
	}

	return &Reader{source: source}, nil
}

func (r *Reader) Read(ctx context.Context, blocks []datanode.BlockInfo, dst io.Writer) (int64, error) {
	var total int64

	for _, block := range blocks {
		rc, err := r.source.GetBlock(ctx, block.ID)
		if err != nil {
			return total, err
		}

		n, copyErr := io.Copy(dst, rc)
		closeErr := rc.Close()

		total += n

		if copyErr != nil {
			return total, copyErr
		}
		if closeErr != nil {
			return total, closeErr
		}
	}

	return total, nil
}
