package slicer

// Chunk divide slice per maxSize
// example []{1,2,3,4,5,6,7} with max 3
// result is []{1,2,3} []{4,5,6} []{7}
func Chunk[T any](xs []T, maxSize int) [][]T {
	if len(xs) == 0 {
		return nil
	}
	divided := make([][]T, (len(xs)+maxSize-1)/maxSize)
	prev := 0
	i := 0
	till := len(xs) - maxSize
	for prev < till {
		next := prev + maxSize
		divided[i] = xs[prev:next]
		prev = next
		i++
	}
	divided[i] = xs[prev:]
	return divided
}
