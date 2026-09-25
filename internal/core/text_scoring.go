package core

import (
	"fmt"
	"math"
)

type TextScoring struct {
	items  [][]float64
	dims   int
	length int
}

func newTextScoring(dims int, text []byte) (*TextScoring, error) {
	if dims == 0 {
		return nil, fmt.Errorf("no evaluations provided")
	}
	length := len(text)
	if length == 0 {
		return nil, fmt.Errorf("empty text")
	}

	var items [][]float64

	for range dims {
		items = append(items, make([]float64, length))
	}

	return &TextScoring{
		items:  items,
		dims:   dims,
		length: length,
	}, nil
}

func (t *TextScoring) Items() [][]float64 {
	return t.items
}

func (t *TextScoring) Length() int {
	return t.length
}

func (t *TextScoring) Dims() int {
	return t.dims
}

func (t *TextScoring) SetAtIndex(dim int, idx int, val float64) error {
	if idx < 0 {
		return fmt.Errorf("idx must be > 0 but has value %d", idx)
	}

	if dim < 0 {
		return fmt.Errorf("dim must be > 0 but has value %d", dim)
	}

	if idx >= t.length {
		return fmt.Errorf("idx must be < %d but is %d", t.length, idx)
	}

	if dim >= t.dims {
		return fmt.Errorf("dim must be < %d but has value %d", t.dims, dim)
	}

	t.items[dim][idx] = val

	return nil
}

func (t *TextScoring) SetInRange(dim, start, end int, val float64) error {
	if dim < 0 {
		return fmt.Errorf("dim must be > 0 but has value %d", dim)
	}

	if dim >= t.dims {
		return fmt.Errorf("dim must be < %d but has value %d", t.dims, dim)
	}

	if start >= end {
		return fmt.Errorf("start (%d) has to be less than end (%d)", start, end)
	}

	if start < 0 || end < 0 {
		return fmt.Errorf("start (%d) and end (%d) must be >= 0", start, end)
	}

	if start >= t.length || end >= t.length {
		return fmt.Errorf("start (%d) and end (%d) must be < length", start, end)
	}

	for idx := start; idx < end; idx++ {
		if err := t.SetAtIndex(dim, idx, val); err != nil {
			return err
		}
	}

	return nil
}

func (t *TextScoring) GetAtIndex(dim, idx int) (float64, error) {
	if dim < 0 {
		return -math.MaxFloat64, fmt.Errorf("dim must be > 0 but has value %d", dim)
	}

	if dim >= t.dims {
		return -math.MaxFloat64, fmt.Errorf("dim must be < %d but has value %d", t.dims, dim)
	}

	if idx < 0 {
		return -math.MaxFloat64, fmt.Errorf("idx must be >= 0 but was %d", idx)
	}

	if idx >= t.length {
		return -math.MaxFloat64, fmt.Errorf("idx must be < %d but was %d", idx, t.length)
	}

	return t.items[dim][idx], nil
}

func (t *TextScoring) GetInRange(dim, start, end int) ([]float64, error) {
	if dim < 0 {
		return nil, fmt.Errorf("dim must be > 0 but has value %d", dim)
	}

	if dim >= t.dims {
		return nil, fmt.Errorf("dim must be < %d but has value %d", t.dims, dim)
	}

	if start >= end {
		return nil, fmt.Errorf("start (%d) has to be less than end (%d)", start, end)
	}

	if start < 0 || end < 0 {
		return nil, fmt.Errorf("start (%d) and end (%d) must be >= 0", start, end)
	}

	if start >= t.length || end >= t.length {
		return nil, fmt.Errorf("start (%d) and end (%d) must be < length", start, end)
	}

	output := make([]float64, end-start)

	for idx := start; idx < end; idx++ {
		output[idx] = t.items[dim][idx]
	}

	return output, nil
}
