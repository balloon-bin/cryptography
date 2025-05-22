package matrix_test

import (
	"encoding/json"
	"testing"

	"git.omicron.one/playground/cryptography/matrix"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	m := matrix.CreateFromFlatSlice(2, 3, []int{1, 2, 3, 4, 5, 6})
	assert.NotNil(t, m)

	mf := matrix.Convert[float64](m)
	assert.Equal(t, m.Rows(), mf.Rows())
	assert.Equal(t, m.Cols(), mf.Cols())

	for row := range mf.Rows() {
		for col := range mf.Cols() {
			assert.Equal(t, float64(m.Get(row, col)), mf.Get(row, col))
		}
	}
}

func TestCreate(t *testing.T) {
	m := matrix.Create[int](3, 4)
	assert.NotNil(t, m)

	rows, cols := m.Size()
	assert.Equal(t, 3, rows)
	assert.Equal(t, 3, m.Rows())
	assert.Equal(t, 4, cols)
	assert.Equal(t, 4, m.Cols())

	for row := range rows {
		for col := range cols {
			assert.Equal(t, 0, m.Get(row, col))
		}
	}

	assert.PanicsWithValue(t, matrix.ErrInvalidDimensions, func() {
		matrix.Create[int](0, 1)
	})
	assert.PanicsWithValue(t, matrix.ErrInvalidDimensions, func() {
		matrix.Create[int](-1, 1)
	})
	assert.PanicsWithValue(t, matrix.ErrInvalidDimensions, func() {
		matrix.Create[int](1, 0)
	})
	assert.PanicsWithValue(t, matrix.ErrInvalidDimensions, func() {
		matrix.Create[int](1, -1)
	})
	assert.PanicsWithValue(t, matrix.ErrInvalidDimensions, func() {
		matrix.Create[int](-1, -1)
	})
}

func TestCreateFromSlice(t *testing.T) {
	m := matrix.CreateFromSlice([][]int{
		{1, 2, 3},
		{4, 5, 6},
	})
	assert.NotNil(t, m)

	assert.Equal(t, 2, m.Rows())
	assert.Equal(t, 3, m.Cols())

	assert.Equal(t, 1, m.Get(0, 0))
	assert.Equal(t, 2, m.Get(0, 1))
	assert.Equal(t, 3, m.Get(0, 2))
	assert.Equal(t, 4, m.Get(1, 0))
	assert.Equal(t, 5, m.Get(1, 1))
	assert.Equal(t, 6, m.Get(1, 2))

	assert.PanicsWithValue(t, matrix.ErrInvalidDimensions, func() {
		matrix.CreateFromSlice([][]int{})
	})

	assert.PanicsWithValue(t, matrix.ErrInvalidDimensions, func() {
		matrix.CreateFromSlice([][]int{
			{},
		})
	})

	assert.PanicsWithValue(t, matrix.ErrIncompatibleDataDimensions, func() {
		matrix.CreateFromSlice([][]int{
			{1, 2, 3},
			{4, 5},
		})
	})
}

func TestCreateFromFlatSlice(t *testing.T) {
	m := matrix.CreateFromFlatSlice(2, 3, []int{1, 2, 3, 4, 5, 6})
	assert.NotNil(t, m)

	assert.Equal(t, 2, m.Rows())
	assert.Equal(t, 3, m.Cols())

	assert.Equal(t, 1, m.Get(0, 0))
	assert.Equal(t, 2, m.Get(0, 1))
	assert.Equal(t, 3, m.Get(0, 2))
	assert.Equal(t, 4, m.Get(1, 0))
	assert.Equal(t, 5, m.Get(1, 1))
	assert.Equal(t, 6, m.Get(1, 2))

	assert.PanicsWithValue(t, matrix.ErrIncompatibleDataDimensions, func() {
		matrix.CreateFromFlatSlice(2, 3, []int{1, 2, 3, 4, 5})
	})

	assert.PanicsWithValue(t, matrix.ErrIncompatibleDataDimensions, func() {
		matrix.CreateFromFlatSlice(2, 3, []int{1, 2, 3, 4, 5, 6, 7})
	})

	assert.PanicsWithValue(t, matrix.ErrInvalidDimensions, func() {
		matrix.CreateFromFlatSlice(0, 1, []int{})
	})
	assert.PanicsWithValue(t, matrix.ErrInvalidDimensions, func() {
		matrix.CreateFromFlatSlice(-1, 1, []int{})
	})
	assert.PanicsWithValue(t, matrix.ErrInvalidDimensions, func() {
		matrix.CreateFromFlatSlice(1, 0, []int{})
	})
	assert.PanicsWithValue(t, matrix.ErrInvalidDimensions, func() {
		matrix.CreateFromFlatSlice(1, -1, []int{})
	})
	assert.PanicsWithValue(t, matrix.ErrInvalidDimensions, func() {
		matrix.CreateFromFlatSlice(-1, -1, []int{1})
	})
}

func TestCreateFromJSON(t *testing.T) {
	// data json
	data := []byte(`[[1, 2, 3], [4, 5, 6]]`)
	m, err := matrix.CreateFromJSON[int](data)
	assert.Nil(t, err)
	assert.NotNil(t, m)
	assert.Equal(t, 2, m.Rows())
	assert.Equal(t, 3, m.Cols())
	assert.Equal(t, 1, m.Get(0, 0))
	assert.Equal(t, 2, m.Get(0, 1))
	assert.Equal(t, 3, m.Get(0, 2))
	assert.Equal(t, 4, m.Get(1, 0))
	assert.Equal(t, 5, m.Get(1, 1))
	assert.Equal(t, 6, m.Get(1, 2))

	// invalid json
	data = []byte(`[[1, 2, 3], [4, 5,`)
	m, err = matrix.CreateFromJSON[int](data)
	assert.NotNil(t, err)
	assert.Nil(t, m)

	// empty matrix
	data = []byte(`[]`)
	m, err = matrix.CreateFromJSON[int](data)
	assert.ErrorIs(t, err, matrix.ErrInvalidDimensions)
	assert.Nil(t, m)

	// empty rows
	data = []byte(`[[]]`)
	m, err = matrix.CreateFromJSON[int](data)
	assert.ErrorIs(t, err, matrix.ErrInvalidDimensions)
	assert.Nil(t, m)

	// mixed row length
	data = []byte(`[[1, 2, 3], [4, 5]]`)
	m, err = matrix.CreateFromJSON[int](data)
	assert.ErrorIs(t, err, matrix.ErrIncompatibleDataDimensions)
	assert.Nil(t, m)
}

func TestSum(t *testing.T) {
	a := matrix.CreateFromFlatSlice(2, 3, []int{1, 2, 3, 4, 5, 6})
	b := matrix.CreateFromFlatSlice(2, 3, []int{1, 1, 1, 1, 1, 1})

	// Sum only one matrix
	c := matrix.Sum(a)
	assert.NotNil(t, c)
	assert.NotSame(t, a, c)

	assert.Equal(t, 2, c.Rows())
	assert.Equal(t, 3, c.Cols())

	assert.Equal(t, 1, c.Get(0, 0))
	assert.Equal(t, 2, c.Get(0, 1))
	assert.Equal(t, 3, c.Get(0, 2))
	assert.Equal(t, 4, c.Get(1, 0))
	assert.Equal(t, 5, c.Get(1, 1))
	assert.Equal(t, 6, c.Get(1, 2))

	// Sum one matrix multiple times
	c = matrix.Sum(a, a, a)
	assert.NotNil(t, c)
	assert.NotSame(t, a, c)

	assert.Equal(t, 2, c.Rows())
	assert.Equal(t, 3, c.Cols())

	assert.Equal(t, 3, c.Get(0, 0))
	assert.Equal(t, 6, c.Get(0, 1))
	assert.Equal(t, 9, c.Get(0, 2))
	assert.Equal(t, 12, c.Get(1, 0))
	assert.Equal(t, 15, c.Get(1, 1))
	assert.Equal(t, 18, c.Get(1, 2))

	// Sum different matrices
	c = matrix.Sum(a, b)
	assert.NotNil(t, c)
	assert.NotEqual(t, a, c)

	assert.Equal(t, 2, c.Rows())
	assert.Equal(t, 3, c.Cols())

	assert.Equal(t, 2, c.Get(0, 0))
	assert.Equal(t, 3, c.Get(0, 1))
	assert.Equal(t, 4, c.Get(0, 2))
	assert.Equal(t, 5, c.Get(1, 0))
	assert.Equal(t, 6, c.Get(1, 1))
	assert.Equal(t, 7, c.Get(1, 2))

	// Sum incorrect dimensions
	d := matrix.Create[int](3, 2)
	assert.PanicsWithValue(t, matrix.ErrIncompatibleMatrixDimensions, func() {
		matrix.Sum(a, d)
	})
}

func TestTransform(t *testing.T) {
	a := matrix.CreateFromFlatSlice(2, 3, []int{1, 2, 3, 4, 5, 6})
	b := matrix.Transform(a, func(value int) complex128 {
		return complex(float64(value), 3.0)
	})

	assert.Equal(t, a.Rows(), b.Rows())
	assert.Equal(t, a.Cols(), b.Cols())

	assert.Equal(t, complex(1.0, 3.0), b.Get(0, 0))
	assert.Equal(t, complex(2.0, 3.0), b.Get(0, 1))
	assert.Equal(t, complex(3.0, 3.0), b.Get(0, 2))
	assert.Equal(t, complex(4.0, 3.0), b.Get(1, 0))
	assert.Equal(t, complex(5.0, 3.0), b.Get(1, 1))
	assert.Equal(t, complex(6.0, 3.0), b.Get(1, 2))
}

func TestMatrix_Add(t *testing.T) {
	a := matrix.CreateFromFlatSlice(2, 3, []int{1, 2, 3, 4, 5, 6})
	b := matrix.CreateFromFlatSlice(2, 3, []int{1, 1, 1, 1, 1, 1})

	// Add nothing
	c := a.Add()
	assert.Same(t, a, c)
	assert.Equal(t, 1, a.Get(0, 0))
	assert.Equal(t, 2, a.Get(0, 1))
	assert.Equal(t, 3, a.Get(0, 2))
	assert.Equal(t, 4, a.Get(1, 0))
	assert.Equal(t, 5, a.Get(1, 1))
	assert.Equal(t, 6, a.Get(1, 2))

	// Add itself multiple times
	c = a.Add(a, a, a)
	assert.Same(t, a, c)

	assert.Equal(t, 4, a.Get(0, 0))
	assert.Equal(t, 8, a.Get(0, 1))
	assert.Equal(t, 12, a.Get(0, 2))
	assert.Equal(t, 16, a.Get(1, 0))
	assert.Equal(t, 20, a.Get(1, 1))
	assert.Equal(t, 24, a.Get(1, 2))

	// Add other matrix
	c = a.Add(b)
	assert.Same(t, a, c)

	assert.Equal(t, 5, a.Get(0, 0))
	assert.Equal(t, 9, a.Get(0, 1))
	assert.Equal(t, 13, a.Get(0, 2))
	assert.Equal(t, 17, a.Get(1, 0))
	assert.Equal(t, 21, a.Get(1, 1))
	assert.Equal(t, 25, a.Get(1, 2))

	// Add incorrect dimension
	assert.PanicsWithValue(
		t, matrix.ErrIncompatibleMatrixDimensions,
		func() {
			d := matrix.Create[int](3, 2)
			a.Add(b, a, a, d)
		},
	)

	assert.Equal(t, 5, a.Get(0, 0))
	assert.Equal(t, 9, a.Get(0, 1))
	assert.Equal(t, 13, a.Get(0, 2))
	assert.Equal(t, 17, a.Get(1, 0))
	assert.Equal(t, 21, a.Get(1, 1))
	assert.Equal(t, 25, a.Get(1, 2))
}

func TestMatrix_Subtract(t *testing.T) {
	a := matrix.CreateFromFlatSlice(2, 3, []int{1, 2, 3, 4, 5, 6})
	b := matrix.CreateFromFlatSlice(2, 3, []int{1, 1, 1, 1, 1, 1})

	// Subtract nothing
	c := a.Subtract()
	assert.Same(t, a, c)
	assert.Equal(t, 1, a.Get(0, 0))
	assert.Equal(t, 2, a.Get(0, 1))
	assert.Equal(t, 3, a.Get(0, 2))
	assert.Equal(t, 4, a.Get(1, 0))
	assert.Equal(t, 5, a.Get(1, 1))
	assert.Equal(t, 6, a.Get(1, 2))

	// Add itself multiple times
	c = a.Subtract(a, a, a, a)
	assert.Same(t, a, c)

	assert.Equal(t, -3, a.Get(0, 0))
	assert.Equal(t, -6, a.Get(0, 1))
	assert.Equal(t, -9, a.Get(0, 2))
	assert.Equal(t, -12, a.Get(1, 0))
	assert.Equal(t, -15, a.Get(1, 1))
	assert.Equal(t, -18, a.Get(1, 2))

	// Add other matrix
	c = a.Subtract(b)
	assert.Same(t, a, c)

	assert.Equal(t, -4, a.Get(0, 0))
	assert.Equal(t, -7, a.Get(0, 1))
	assert.Equal(t, -10, a.Get(0, 2))
	assert.Equal(t, -13, a.Get(1, 0))
	assert.Equal(t, -16, a.Get(1, 1))
	assert.Equal(t, -19, a.Get(1, 2))

	// Add incorrect dimension
	assert.PanicsWithValue(
		t, matrix.ErrIncompatibleMatrixDimensions,
		func() {
			d := matrix.Create[int](3, 2)
			a.Subtract(b, a, a, d)
		},
	)

	assert.Equal(t, -4, a.Get(0, 0))
	assert.Equal(t, -7, a.Get(0, 1))
	assert.Equal(t, -10, a.Get(0, 2))
	assert.Equal(t, -13, a.Get(1, 0))
	assert.Equal(t, -16, a.Get(1, 1))
	assert.Equal(t, -19, a.Get(1, 2))
}

func TestMatrix_Apply(t *testing.T) {
	a := matrix.CreateFromFlatSlice(2, 3, []int{1, 2, 3, 4, 5, 6})
	assert.NotNil(t, a)

	a.Apply(func(x int) int { return x * x })

	assert.Equal(t, 1, a.Get(0, 0))
	assert.Equal(t, 4, a.Get(0, 1))
	assert.Equal(t, 9, a.Get(0, 2))
	assert.Equal(t, 16, a.Get(1, 0))
	assert.Equal(t, 25, a.Get(1, 1))
	assert.Equal(t, 36, a.Get(1, 2))
}

func TestMatrix_Copy(t *testing.T) {
	a := matrix.CreateFromFlatSlice(2, 3, []int{1, 2, 3, 4, 5, 6})
	assert.NotNil(t, a)

	b := a.Copy()
	assert.NotSame(t, a, b)
	assert.Equal(t, 1, b.Get(0, 0))
	assert.Equal(t, 2, b.Get(0, 1))
	assert.Equal(t, 3, b.Get(0, 2))
	assert.Equal(t, 4, b.Get(1, 0))
	assert.Equal(t, 5, b.Get(1, 1))
	assert.Equal(t, 6, b.Get(1, 2))
}

func TestMatrix_Get(t *testing.T) {
	a := matrix.CreateFromFlatSlice(2, 3, []int{1, 2, 3, 4, 5, 6})
	assert.NotNil(t, a)

	assert.Equal(t, 1, a.Get(0, 0))
	assert.Equal(t, 2, a.Get(0, 1))
	assert.Equal(t, 3, a.Get(0, 2))
	assert.Equal(t, 4, a.Get(1, 0))
	assert.Equal(t, 5, a.Get(1, 1))
	assert.Equal(t, 6, a.Get(1, 2))

	assert.Panics(t, func() {
		a.Get(-1, 0)
	})
	assert.Panics(t, func() {
		a.Get(0, -1)
	})
	assert.Panics(t, func() {
		a.Get(2, 0)
	})
	assert.Panics(t, func() {
		a.Get(0, 3)
	})
}

func TestMatrix_Set(t *testing.T) {
	a := matrix.Create[int](2, 2)
	assert.NotNil(t, a)

	for row := range a.Rows() {
		for col := range a.Cols() {
			assert.NotEqual(t, a.Get(row, col), 42)
			a.Set(row, col, 42)
			assert.Equal(t, 42, a.Get(row, col))
		}
	}

	assert.Panics(t, func() {
		a.Set(-1, 0, 42)
	})
	assert.Panics(t, func() {
		a.Set(0, -1, 42)
	})
	assert.Panics(t, func() {
		a.Set(2, 0, 42)
	})
	assert.Panics(t, func() {
		a.Set(0, 2, 42)
	})
}

func TestMatrix_Scale(t *testing.T) {
	a := matrix.CreateFromFlatSlice(2, 3, []int{1, 2, 3, 4, 5, 6})

	// Add nothing
	b := a.Scale(-4)
	assert.Same(t, a, b)

	assert.Equal(t, -4, a.Get(0, 0))
	assert.Equal(t, -8, a.Get(0, 1))
	assert.Equal(t, -12, a.Get(0, 2))
	assert.Equal(t, -16, a.Get(1, 0))
	assert.Equal(t, -20, a.Get(1, 1))
	assert.Equal(t, -24, a.Get(1, 2))
}

func TestMatrix_HadamardMultiply(t *testing.T) {
	a := matrix.CreateFromFlatSlice(2, 3, []int{1, 2, 3, 4, 5, 6})
	b := matrix.CreateFromFlatSlice(2, 3, []int{2, 2, 2, 2, 2, 2})

	// Multiply nothing
	c := a.HadamardMultiply()
	assert.Same(t, a, c)
	assert.Equal(t, 1, a.Get(0, 0))
	assert.Equal(t, 2, a.Get(0, 1))
	assert.Equal(t, 3, a.Get(0, 2))
	assert.Equal(t, 4, a.Get(1, 0))
	assert.Equal(t, 5, a.Get(1, 1))
	assert.Equal(t, 6, a.Get(1, 2))

	// Multiply itself multiple times
	c = a.HadamardMultiply(a, a, a)
	assert.Same(t, a, c)

	assert.Equal(t, 1, a.Get(0, 0))
	assert.Equal(t, 16, a.Get(0, 1))
	assert.Equal(t, 81, a.Get(0, 2))
	assert.Equal(t, 256, a.Get(1, 0))
	assert.Equal(t, 625, a.Get(1, 1))
	assert.Equal(t, 1296, a.Get(1, 2))

	// Multiply other matrix
	c = a.HadamardMultiply(b)
	assert.Same(t, a, c)

	assert.Equal(t, 2, a.Get(0, 0))
	assert.Equal(t, 32, a.Get(0, 1))
	assert.Equal(t, 162, a.Get(0, 2))
	assert.Equal(t, 512, a.Get(1, 0))
	assert.Equal(t, 1250, a.Get(1, 1))
	assert.Equal(t, 2592, a.Get(1, 2))

	// Multiply incorrect dimension
	assert.PanicsWithValue(
		t, matrix.ErrIncompatibleMatrixDimensions,
		func() {
			d := matrix.Create[int](3, 2)
			a.HadamardMultiply(b, a, a, d)
		},
	)

	assert.Equal(t, 2, a.Get(0, 0))
	assert.Equal(t, 32, a.Get(0, 1))
	assert.Equal(t, 162, a.Get(0, 2))
	assert.Equal(t, 512, a.Get(1, 0))
	assert.Equal(t, 1250, a.Get(1, 1))
	assert.Equal(t, 2592, a.Get(1, 2))
}

func TestHadamardProduct(t *testing.T) {
	a := matrix.CreateFromFlatSlice(2, 3, []int{1, 2, 3, 4, 5, 6})
	b := matrix.CreateFromFlatSlice(2, 3, []int{2, 2, 2, 2, 2, 2})

	// Multiply only one matrix
	c := matrix.HadamardProduct(a)
	assert.NotNil(t, c)
	assert.NotSame(t, a, c)

	assert.Equal(t, 2, c.Rows())
	assert.Equal(t, 3, c.Cols())

	assert.Equal(t, 1, c.Get(0, 0))
	assert.Equal(t, 2, c.Get(0, 1))
	assert.Equal(t, 3, c.Get(0, 2))
	assert.Equal(t, 4, c.Get(1, 0))
	assert.Equal(t, 5, c.Get(1, 1))
	assert.Equal(t, 6, c.Get(1, 2))

	// Multilply one matrix multiple times
	c = matrix.HadamardProduct(a, a, a)
	assert.NotNil(t, c)
	assert.NotSame(t, a, c)

	assert.Equal(t, 2, c.Rows())
	assert.Equal(t, 3, c.Cols())

	assert.Equal(t, 1, c.Get(0, 0))
	assert.Equal(t, 8, c.Get(0, 1))
	assert.Equal(t, 27, c.Get(0, 2))
	assert.Equal(t, 64, c.Get(1, 0))
	assert.Equal(t, 125, c.Get(1, 1))
	assert.Equal(t, 216, c.Get(1, 2))

	// Multiply different matrices
	c = matrix.HadamardProduct(a, b)
	assert.NotNil(t, c)
	assert.NotEqual(t, a, c)

	assert.Equal(t, 2, c.Rows())
	assert.Equal(t, 3, c.Cols())

	assert.Equal(t, 2, c.Get(0, 0))
	assert.Equal(t, 4, c.Get(0, 1))
	assert.Equal(t, 6, c.Get(0, 2))
	assert.Equal(t, 8, c.Get(1, 0))
	assert.Equal(t, 10, c.Get(1, 1))
	assert.Equal(t, 12, c.Get(1, 2))

	// Multiply incorrect dimensions
	d := matrix.Create[int](3, 2)
	assert.PanicsWithValue(t, matrix.ErrIncompatibleMatrixDimensions, func() {
		matrix.HadamardProduct(a, d)
	})
}

func TestMatrix_Fill(t *testing.T) {
	a := matrix.Create[int](2, 3)
	assert.Equal(t, 0, a.Get(0, 0))
	assert.Equal(t, 0, a.Get(0, 1))
	assert.Equal(t, 0, a.Get(0, 2))
	assert.Equal(t, 0, a.Get(1, 0))
	assert.Equal(t, 0, a.Get(1, 1))
	assert.Equal(t, 0, a.Get(1, 2))

	a.Fill(3)
	assert.Equal(t, 3, a.Get(0, 0))
	assert.Equal(t, 3, a.Get(0, 1))
	assert.Equal(t, 3, a.Get(0, 2))
	assert.Equal(t, 3, a.Get(1, 0))
	assert.Equal(t, 3, a.Get(1, 1))
	assert.Equal(t, 3, a.Get(1, 2))
}

func TestMatrix_UnmarshalJSON(t *testing.T) {
	// int matrix
	data := []byte(`[[1,2,3],[4,5,6]]`)
	var m *matrix.Matrix[int]
	err := json.Unmarshal(data, &m)
	assert.Nil(t, err)

	assert.Equal(t, 2, m.Rows())
	assert.Equal(t, 3, m.Cols())
	assert.Equal(t, 1, m.Get(0, 0))
	assert.Equal(t, 2, m.Get(0, 1))
	assert.Equal(t, 3, m.Get(0, 2))
	assert.Equal(t, 4, m.Get(1, 0))
	assert.Equal(t, 5, m.Get(1, 1))
	assert.Equal(t, 6, m.Get(1, 2))

	// float matrix
	data = []byte(`[[1.5,2.5],[3.5,4.5]]`)
	var mf *matrix.Matrix[float64]
	err = json.Unmarshal(data, &mf)
	assert.Nil(t, err)

	assert.Equal(t, 2, mf.Rows())
	assert.Equal(t, 2, mf.Cols())
	assert.Equal(t, 1.5, mf.Get(0, 0))
	assert.Equal(t, 2.5, mf.Get(0, 1))
	assert.Equal(t, 3.5, mf.Get(1, 0))
	assert.Equal(t, 4.5, mf.Get(1, 1))

	// via json.Unmarshal
	matrices := []byte(`[[[1,2],[3,4]],[[5,6,7]]]`)
	var ms []*matrix.Matrix[int]
	err = json.Unmarshal(matrices, &ms)
	assert.Nil(t, err)
	assert.Len(t, ms, 2)
	assert.Equal(t, 2, ms[0].Get(0, 1))
	assert.Equal(t, 7, ms[1].Get(0, 2))

	// invalid JSON
	err = m.UnmarshalJSON([]byte(`invalid`))
	assert.NotNil(t, err)

	// empty array
	err = m.UnmarshalJSON([]byte(`[]`))
	assert.ErrorIs(t, err, matrix.ErrInvalidDimensions)

	// empty inner array
	err = m.UnmarshalJSON([]byte(`[[]]`))
	assert.ErrorIs(t, err, matrix.ErrInvalidDimensions)

	// inconsistent lengths
	err = m.UnmarshalJSON([]byte(`[[1,2],[3]]`))
	assert.ErrorIs(t, err, matrix.ErrIncompatibleDataDimensions)
}

func TestMatrix_MarshallJSON(t *testing.T) {
	// int matrix
	m := matrix.CreateFromSlice([][]int{{1, 2, 3}, {4, 5, 6}})
	data, err := m.MarshalJSON()
	assert.Nil(t, err)
	assert.NotNil(t, data)

	expected := `[[1,2,3],[4,5,6]]`
	assert.Equal(t, expected, string(data))

	// float matrix
	mf := matrix.CreateFromSlice([][]float64{{1.5, 2.5}, {3.5, 4.5}})
	data, err = mf.MarshalJSON()
	assert.Nil(t, err)
	assert.NotNil(t, data)
	expectedFloat := `[[1.5,2.5],[3.5,4.5]]`
	assert.Equal(t, expectedFloat, string(data))

	// slice of matrices via json.Marshal
	m1 := matrix.CreateFromSlice([][]int{{1, 2}, {3, 4}})
	m2 := matrix.CreateFromSlice([][]int{{5, 6, 7}})
	matrices := []*matrix.Matrix[int]{m1, m2}

	data, err = json.Marshal(matrices)
	assert.Nil(t, err)
	expected = `[[[1,2],[3,4]],[[5,6,7]]]`
	assert.Equal(t, expected, string(data))
}
