// matrix provides very basic matrix operations with all numeric types. It is
// not a high performance or feature rich implementation and is mostly intended
// to use for collecting statistics
package matrix

import (
	"errors"

	"golang.org/x/exp/constraints"
)

var (
	ErrInvalidDimensions            = errors.New("Invalid dimensions")
	ErrIncompatibleMatrixDimensions = errors.New("Incompatible matrix dimensions")
	ErrIncompatibleDataDimensions   = errors.New("Incompatible data dimensions")
)

// Matrices can be created for these underlying types
type Number interface {
	constraints.Integer | constraints.Float | constraints.Complex
}

type SimpleNumber interface {
	constraints.Integer | constraints.Float
}

// Matrix represents a matrix with values of a specific Number type
type Matrix[T Number] struct {
	rows   int
	cols   int
	values [][]T
}

func createValues[T Number](rows, cols int) [][]T {
	values := make([][]T, rows)
	for i := range rows {
		values[i] = make([]T, cols)
	}
	return values
}

// Create creates a new matrix with the given number of rows and columns. All
// values of this matrix are set to zero. Returns the new matrix.
//
// Panics with ErrInvalidDimensions if rows < 1 or cols < 1.
func Create[T Number](rows, cols int) *Matrix[T] {
	if rows < 1 || cols < 1 {
		panic(ErrInvalidDimensions)
	}
	return &Matrix[T]{
		rows:   rows,
		cols:   cols,
		values: createValues[T](rows, cols),
	}
}

// Creates a new matrix of a given size and sets the values from the given
// slice. The values are used to fill the matrix left to right and top to
// bottom. Returns the newly created matrix.
//
// Panics with ErrInvalidDimensions if rows < 1 or cols < 1.
// Panics with ErrIncompatibleDataDimensions if the length of the values doesn't
// match the size of the matrix. This is the only possible error condition.
func CreateFromFlatSlice[T Number](rows, cols int, values []T) *Matrix[T] {
	if rows < 1 || cols < 1 {
		panic(ErrInvalidDimensions)
	}
	if len(values) != rows*cols {
		panic(ErrIncompatibleDataDimensions)
	}
	m := Create[T](rows, cols)

	n := 0
	for i := range rows {
		for j := range cols {
			m.values[i][j] = values[n]
			n++
		}
	}
	return m
}

// Convert will take a matrix of type T and convert it into a matrix of type U
// Only works on SimpleNumber matrices
func Convert[U, T SimpleNumber](in *Matrix[T]) *Matrix[U] {
	out := Create[U](in.rows, in.cols)
	for i := range in.rows {
		for j := range in.cols {
			out.values[i][j] = U(in.values[i][j])
		}
	}
	return out
}

// Transform the values of the given matrix with a transfrom function. This
// operation may change the type of the values. Unlike Convert it can be used
// on all Number types, not just SimpleNumber.
func Transform[U, T Number](in *Matrix[T], transformFn func(T) U) *Matrix[U] {
	out := Create[U](in.rows, in.cols)
	for i := range in.rows {
		for j := range in.cols {
			out.values[i][j] = transformFn(in.values[i][j])
		}
	}
	return out
}

// Sum takes at least one matrix and sums it together.
// Returns a new matrix that is the sum of the arguments.
// Panics if the arguments don't have matching dimensions.
func Sum[T Number](first *Matrix[T], additional ...*Matrix[T]) *Matrix[T] {
	return first.Copy().Add(additional...)
}

// HadamardProduct takes at least one matrix and performs component-wise
// multiplication of all given matrices. Returns a new matrix with the computed
// values. Panics if the arguments don't have matching dimensions.
func HadamardProduct[T Number](first *Matrix[T], additional ...*Matrix[T]) *Matrix[T] {
	return first.Copy().HadamardMultiply(additional...)
}

// Copy creates a deep copy of the matrix and returns the new instance
func (m *Matrix[T]) Copy() *Matrix[T] {
	mCopy := Create[T](m.rows, m.cols)
	for i := range m.rows {
		copy(mCopy.values[i], m.values[i])
	}
	return mCopy
}

// Size returns the dimensions of the matrix as (rows, columns)
func (m *Matrix[T]) Size() (int, int) {
	return m.rows, m.cols
}

// Rows returns the number of rows in the matrix
func (m *Matrix[T]) Rows() int {
	return m.rows
}

// Cols returns the number of columns in the matrix
func (m *Matrix[T]) Cols() int {
	return m.cols
}

// Set sets the value of the matrix at the given row and col position to the
// given value
func (m *Matrix[T]) Set(row, col int, value T) {
	m.values[row][col] = value
}

// Set assigns the specified value to the element at the given row and column
func (m *Matrix[T]) Get(row, col int) T {
	return m.values[row][col]
}

// Add performs in-place addition of zero or more matrices to this matrix and
// returns the receiver.
// Ensures correct behavior even if the matrix itself is passed as one or more
// arguments.
// Panics with ErrIncompatibleMatrixDimensions if any of the matrices don't
// have matching dimensions.
func (m *Matrix[T]) Add(matrices ...*Matrix[T]) *Matrix[T] {
	numSelf := 0
	for _, other := range matrices {
		if m.rows != other.rows || m.cols != other.cols {
			panic(ErrIncompatibleMatrixDimensions)
		}
		if other == m {
			numSelf += 1
		}
	}

	// If we have multiple self references to add, work on duplicate data to
	// make addition behave as expected
	values := m.values
	if numSelf > 1 {
		values = m.Copy().values
	}

	for _, other := range matrices {
		for i := range m.rows {
			for j := range m.cols {
				values[i][j] += other.values[i][j]
			}
		}
	}

	m.values = values
	return m
}

// Subtract performs in-place subtraction of zero or more matrices from this matrix and
// returns the receiver.
// Ensures correct behavior even if the matrix itself is passed as one or more
// arguments.
// Panics with ErrIncompatibleMatrixDimensions if any of the matrices don't
// have matching dimensions.
func (m *Matrix[T]) Subtract(matrices ...*Matrix[T]) *Matrix[T] {
	numSelf := 0
	for _, other := range matrices {
		if m.rows != other.rows || m.cols != other.cols {
			panic(ErrIncompatibleMatrixDimensions)
		}
		if other == m {
			numSelf += 1
		}
	}

	// If we have multiple self references to subtract, work on duplicate data to
	// make subtraction behave as expected
	values := m.values
	if numSelf > 1 {
		values = m.Copy().values
	}

	for _, other := range matrices {
		for i := range m.rows {
			for j := range m.cols {
				values[i][j] -= other.values[i][j]
			}
		}
	}

	m.values = values
	return m
}

// Apply performs an in-place transformation of each element of the matrix using
// the provided function. Returns the receiver.
func (m *Matrix[T]) Apply(fn func(T) T) *Matrix[T] {
	for i := range m.rows {
		for j := range m.cols {
			m.values[i][j] = fn(m.values[i][j])
		}
	}
	return m
}

// Scale does an in-place scalar multiplication of the matrix values. Returns
// the receiver.
func (m *Matrix[T]) Scale(scalar T) *Matrix[T] {
	for i := range m.rows {
		for j := range m.cols {
			m.values[i][j] *= scalar
		}
	}
	return m
}

// HadamardMultiply performs an in-place component-wise multiplication of this
// matrix with zero or more matrices.
// Ensures correct behavior even if the matrix itself is passed as one or more
// arguments.
// Panics with ErrIncompatibleMatrixDimensions if any of the matrices don't
// have matching dimensions.
func (m *Matrix[T]) HadamardMultiply(matrices ...*Matrix[T]) *Matrix[T] {
	numSelf := 0
	for _, other := range matrices {
		if m.rows != other.rows || m.cols != other.cols {
			panic(ErrIncompatibleMatrixDimensions)
		}
		if other == m {
			numSelf += 1
		}
	}

	// If we have multiple self references to multiply, work on duplicate data to
	// make multiplication behave as expected
	values := m.values
	if numSelf > 1 {
		values = m.Copy().values
	}

	for _, other := range matrices {
		for i := range m.rows {
			for j := range m.cols {
				values[i][j] *= other.values[i][j]
			}
		}
	}

	m.values = values
	return m
}

// Fill sets all components of this matrix to the given value. Returns the receiver.
func (m *Matrix[T]) Fill(value T) *Matrix[T] {
	for i := range m.rows {
		for j := range m.cols {
			m.values[i][j] = value
		}
	}
	return m
}
