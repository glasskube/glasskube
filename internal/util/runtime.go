package util

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	} else {
		return v
	}
}
