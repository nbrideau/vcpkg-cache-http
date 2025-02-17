package main_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	main "github.com/lesomnus/vcpkg-cache-http"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var DescriptionFoo = main.Description{
    Triplet: "x64-linux",
	Name:    "foo",
	Version: "bar",
	Hash:    "baz",
}

type StoreSetup interface {
	New(t *testing.T) (main.Store, error)
}

type StoreTestSuite struct {
	suite.Suite
	Store StoreSetup

	require *require.Assertions
	store   main.Store
}

func (s *StoreTestSuite) SetupTest() {
	s.require = require.New(s.T())

	store, err := s.Store.New(s.T())
	s.require.NoError(err)
	s.store = store
}

func (s *StoreTestSuite) TestAll() {
	ctx := context.Background()

	_, err := s.store.Head(ctx, DescriptionFoo)
	s.require.ErrorIs(err, main.ErrNotExist)

	err = s.store.Get(ctx, DescriptionFoo, io.Discard)
	s.require.ErrorIs(err, main.ErrNotExist)

	data := randomData(s.T())
	err = s.store.Put(ctx, DescriptionFoo, bytes.NewReader(data))
	s.require.NoError(err)

	size, err := s.store.Head(ctx, DescriptionFoo)
	s.require.NoError(err)
	s.require.Equal(len(data), size)

	var received bytes.Buffer
	err = s.store.Get(ctx, DescriptionFoo, &received)
	s.require.NoError(err)
	s.require.Equal(data, received.Bytes())
}

func (s *StoreTestSuite) TestGetNotExist() {
	ctx := context.Background()

	var received bytes.Buffer
	err := s.store.Get(ctx, DescriptionFoo, &received)
	s.require.ErrorIs(err, main.ErrNotExist)
}

func (s *StoreTestSuite) TestPutAlreadyExist() {
	ctx := context.Background()

	err := s.store.Put(ctx, DescriptionFoo, bytes.NewReader(randomData(s.T())))
	s.require.NoError(err)

	err = s.store.Put(ctx, DescriptionFoo, bytes.NewReader(randomData(s.T())))
	s.require.ErrorIs(err, main.ErrExist)
}
