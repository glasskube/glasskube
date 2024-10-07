package suspend

type Result string

const (
	Suspended Result = "suspended"
	Resumed   Result = "resumed"
	UpToDate  Result = "up-to-date"
	Unknown   Result = "unknown"
)
