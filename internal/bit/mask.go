package bit

import "golang.org/x/exp/constraints"

func Has[T constraints.Unsigned](a, b T) bool {
	return a&b != 0
}
