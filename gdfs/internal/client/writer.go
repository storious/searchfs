package client

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"gdfs/internal/datanode"
)

type BlockWriter interface {
	PutBlock(ctx context.Context, id datanode.BlockID, r io.Reader) (datanode.BlockInfo, error)
}

type WriteResult struct {
	Blocks []datanode.BlockInfo
	Size   int64
}

type Writer struct {
	blockSize int64
	target    BlockWriter
}

func NewWriter(blockSize int64, target BlockWriter) (*Writer, error) {
	if blockSize <= 0 {
		return nil, fmt.Errorf("invalid block size: %d", blockSize)
	}
	if target == nil {
		return nil, fmt.Errorf("nil block writer")
	}

	return &Writer{
		blockSize: blockSize,
		target:    target,
	}, nil
}

func (w *Writer) Write(ctx context.Context, r io.Reader) (WriteResult, error) {
	var result WriteResult

	buf := make([]byte, w.blockSize)
	for index := 0; ; index++ {
		n, err := io.ReadFull(r, buf)
		if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
			return WriteResult{}, err
		}
		if n == 0 && err == io.EOF {
			break
		}

		data := make([]byte, n)
		copy(data, buf[:n])

		blockID := makeBlockID(index, data)

		info, err := w.target.PutBlock(ctx, blockID, bytesReader(data))
		if err != nil {
			return WriteResult{}, err
		}

		result.Blocks = append(result.Blocks, info)
		result.Size += info.Size

		if err == io.ErrUnexpectedEOF || err == io.EOF {
			break
		}
	}

	return result, nil
}

func makeBlockID(index int, data []byte) datanode.BlockID {
	sum := sha256.Sum256(data)
	return datanode.BlockID(fmt.Sprintf("block-%06d-%s", index, hex.EncodeToString(sum[:])[:16]))
}

func bytesReader(b []byte) io.Reader {
	return bytes.NewReader(b)
}
