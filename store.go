package main

import (
	"context"
	"fmt"
	"io"
)

type Description struct {
    Triplet string
	Name    string
	Version string
	Hash    string
}

func (d *Description) String() string {
	return fmt.Sprintf("/%s/%s/%s/%s", d.Triplet, d.Name, d.Version, d.Hash)
}

type Store interface {
	Get(ctx context.Context, desc Description, w io.Writer) error
	Head(ctx context.Context, desc Description) (int, error)
	Put(ctx context.Context, desc Description, r io.Reader) error

	Close() error
}
