// Package rpc implements the rpc layer for the bloomfilter, following the principles from https://golang.org/pkg/net/rpc
package rpc

import (
	"context"
	"fmt"

	"api-gateway/v2/modules/bloomfilter/v2/rotate"
)

// Config type containing a sliding bloomfilter set and a port
type Config struct {
	rotate.Config
	Port int `json:"port"`
}

// BloomfilterRPC type
type BloomfilterRPC int

// Bloomfilter wrapper for BloomfilterRPC type
type Bloomfilter struct {
	BloomfilterRPC
}

// New rpc layer implementation of creating a sliding bloomfilter set
func New(ctx context.Context, cfg Config) *Bloomfilter {
	if bf != nil {
		bf.Close()
	}

	bf = rotate.New(ctx, cfg.Config)

	return new(Bloomfilter)
}

// AddInput type for an array of elements
type AddInput struct {
	Elems [][]byte
}

// AddOutput type for an array of elements to a sliding bloomfilter set
type AddOutput struct {
	Count int
}

// Add rpc layer implementation of an array of elements to a sliding bloomfilter set
func (BloomfilterRPC) Add(in AddInput, out *AddOutput) error {
	if bf == nil {
		out.Count = 0
		return ErrNoBloomfilterInitialized
	}

	k := 0
	for _, elem := range in.Elems {
		bf.Add(elem)
		k++
	}
	out.Count = k

	return nil
}

// CheckInput type for an array of elements
type CheckInput struct {
	Elems [][]byte
}

// CheckOutput type for check result of an array of elements in a sliding bloomfilter set
type CheckOutput struct {
	Checks []bool
}

// Check rpc layer implementation of an array of elements in a sliding bloomfilter set
func (BloomfilterRPC) Check(in CheckInput, out *CheckOutput) error {
	checkRes := make([]bool, len(in.Elems))

	if bf == nil {
		out.Checks = checkRes
		return ErrNoBloomfilterInitialized
	}

	for i, elem := range in.Elems {
		checkRes[i] = bf.Check(elem)
	}
	out.Checks = checkRes

	return nil
}

// UnionInput type for sliding bloomfilter set
type UnionInput struct {
	BF *rotate.Bloomfilter
}

// UnionOutput type for sliding bloomfilter set fill degree
type UnionOutput struct {
	Capacity float64
}

// Union rpc layer implementation of two sliding bloomfilter sets
func (BloomfilterRPC) Union(in UnionInput, out *UnionOutput) error {
	if bf == nil {
		out.Capacity = 0
		return ErrNoBloomfilterInitialized
	}

	var err error
	out.Capacity, err = bf.Union(in.BF)

	return err
}

// Close sliding bloomfilter set
func (Bloomfilter) Close() {
	if bf != nil {
		bf.Close()
	}
}

// Bloomfilter getter
func (Bloomfilter) Bloomfilter() *rotate.Bloomfilter {
	return bf
}

// ErrNoBloomfilterInitialized error
var (
	ErrNoBloomfilterInitialized = fmt.Errorf("Bloomfilter not initialized")
	bf                          *rotate.Bloomfilter
)
