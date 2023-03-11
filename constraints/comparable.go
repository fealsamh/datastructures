package constraints

// Comparable is a totally orderable type.
type Comparable[T any] interface {
	Compare(T) int
}
