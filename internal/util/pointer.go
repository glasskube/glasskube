package util

func Pointer[T any](obj T) *T {
	return &obj
}
