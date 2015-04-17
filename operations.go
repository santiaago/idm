package main

import (
	"fmt"
	"math"
)

// times performs a 'a' + 'b' operation and returns it.
func add(a, b Value) Value {
	if _, ok := a.(Int); ok {
		return Int(a.(Int) + b.(Int))
	}
	if _, ok := a.(Vector); ok {
		var v Vector
		for i := 0; i < len(a.(Vector)); i++ {
			v = append(v, add(a.(Vector)[i], b.(Vector)[i]))
		}
		return v
	}
	fmt.Println("ERROR add: case not supported")
	return nil
}

// times performs a 'a' - 'b' operation and returns it.
func minus(a, b Value) Value {
	if _, ok := a.(Int); ok {
		return Int(a.(Int) - b.(Int))
	}
	if _, ok := a.(Vector); ok {
		var v Vector
		for i := 0; i < len(a.(Vector)); i++ {
			v = append(v, minus(a.(Vector)[i], b.(Vector)[i]))
		}
		return v
	}
	fmt.Println("ERROR minus: case not supported")
	return nil
}

// times performs a 'a' * 'b' operation and returns it.
func times(a, b Value) Value {
	if _, ok := a.(Int); ok {
		return Int(a.(Int) * b.(Int))
	}
	if _, ok := a.(Vector); ok {
		var v Vector
		for i := 0; i < len(a.(Vector)); i++ {
			v = append(v, times(a.(Vector)[i], b.(Vector)[i]))
		}
		return v
	}
	fmt.Println("ERROR times: case not supported")
	return nil
}

// pow performs a 'a' ** 'b' operation and returns it.
func pow(a, b Value) Value {
	if _, ok := a.(Int); ok {
		return Int(math.Pow(float64(a.(Int)), float64(b.(Int))))
	}
	if _, ok := a.(Vector); ok {
		var v Vector
		for i := 0; i < len(a.(Vector)); i++ {
			v = append(v, pow(a.(Vector)[i], b.(Vector)[i]))
		}
		return v
	}
	fmt.Println("ERROR pow: case not supported")
	return nil
}

// max performs the maximum value between 'a' and 'b' and returns it.
func max(a, b Value) Value {
	if _, ok := a.(Int); ok {
		return Int(math.Max(float64(a.(Int)), float64(b.(Int))))
	}
	if _, ok := a.(Vector); ok {
		var v Vector
		for i := 0; i < len(a.(Vector)); i++ {
			v = append(v, max(a.(Vector)[i], b.(Vector)[i]))
		}
		return v
	}
	fmt.Println("ERROR max: case not supported")
	return nil
}

// min performs the minimum value between 'a' and 'b' and returns it.
func min(a, b Value) Value {
	if _, ok := a.(Int); ok {
		return Int(math.Min(float64(a.(Int)), float64(b.(Int))))
	}
	if _, ok := a.(Vector); ok {
		var v Vector
		for i := 0; i < len(a.(Vector)); i++ {
			v = append(v, min(a.(Vector)[i], b.(Vector)[i]))
		}
		return v
	}
	fmt.Println("ERROR min: case not supported")
	return nil
}

// sum performs the sum of all items of 'a'. <+/>
// if 'a' is a vector, it is the sum of the vector items.
func sum(a Value) Value {
	if _, ok := a.(Int); ok {
		return a.(Int)
	}
	if _, ok := a.(Vector); ok {
		var v Value
		v = Int(0)
		for i := 0; i < len(a.(Vector)); i++ {
			v = add(v, a.(Vector)[i])
		}
		return v
	}
	fmt.Println("ERROR sum: case not supported")
	return nil
}

// scanSum performs the scan sum of the all the items of 'a'. <+\>
// if 'a' is a vector, the result of scanSum is a vector with the
// cumulative sum of the previous items.
// example +\ 1 2 3
// 1 3 6
func scanSum(a Value) Value {
	if _, ok := a.(Int); ok {
		return a.(Int)
	}

	if _, ok := a.(Vector); ok {
		var v Vector
		for i := 1; i <= len(a.(Vector)); i++ {
			v = append(v, sum(a.(Vector)[:i]))
		}
		return v
	}
	return nil
}

// multiply performs the multiplication of all items of 'a'. <*/>
// if 'a' is a vector, it is the multiplication of the vector items.
func multiply(a Value) Value {
	if _, ok := a.(Int); ok {
		return a.(Int)
	}
	if _, ok := a.(Vector); ok {
		var v Value
		v = Int(1)
		for i := 0; i < len(a.(Vector)); i++ {
			v = times(v, a.(Vector)[i])
		}
		return v
	}
	fmt.Println("ERROR multiply: case not supported")
	return nil
}

// scanMultiply performs the scan multiplication of the all the items of 'a'. <*\>
// if 'a' is a vector, the result of scanSum is a vector with the
// cumulative sum of the previous items.
// example *\ 1 2 3
// 1 2 6
func scanMultiply(a Value) Value {
	if _, ok := a.(Int); ok {
		return a.(Int)
	}

	if _, ok := a.(Vector); ok {
		var v Vector
		for i := 1; i <= len(a.(Vector)); i++ {
			v = append(v, multiply(a.(Vector)[:i]))
		}
		return v
	}
	return nil
}
