package ctrlpkg

func IsSameResource(a, b Package) bool {
	return a.GetName() == b.GetName() &&
		a.GroupVersionKind() == b.GroupVersionKind() &&
		a.GetNamespace() == b.GetNamespace()
}
