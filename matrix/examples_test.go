package matrix_test

import (
	"fmt"

	"git.omicron.one/playground/cryptography/matrix"
)

func ExampleConvert() {
	intMatrix := matrix.Create[int](3, 4)

	_ = matrix.Convert[float64](intMatrix)
}

func ExampleCreateFromFlatSlice() {
	m := matrix.CreateFromFlatSlice(2, 2, []int{0, 1, 2, 3})

	for i := range m.Rows() {
		for j := range m.Cols() {
			fmt.Printf(" %d", m.Get(i, j))
		}
		fmt.Println()
	}

	// Output:
	//  0 1
	//  2 3
}

func ExampleTransform() {
	intMatrix := matrix.CreateFromFlatSlice(2, 2, []int{0, 1, 2, 3})

	_ = matrix.Transform(intMatrix, func(in int) complex128 {
		return complex(float64(in), 0.0)
	})
}

func ExampleMatrix_Apply() {
	m := matrix.CreateFromFlatSlice(2, 2, []int{0, 1, 2, 3})

	m.Apply(func(in int) int {
		return in * in
	})
}
