package condition

type (
	Type   string
	Reason string
)

const (
	Ready  Type = "Ready"
	Failed Type = "Failed"
)

const (
	SyncCompleted             Reason = "SyncCompleted"
	SyncFailed                Reason = "SyncFailed"
	Reconciling               Reason = "Reconciling"
	UpToDate                  Reason = "UpToDate"
	UnsupportedFormat         Reason = "UnsupportedFormat"
	ValueConfigurationInvalid Reason = "ValueConfigurationInvalid"
	InstallationSucceeded     Reason = "InstallationSucceeded"
	InstallationFailed        Reason = "InstallationFailed"
	Pending                   Reason = "Pending"
)

func (r Reason) Recoverable() bool {
	switch r {
	case InstallationFailed:
		return false
	case UnsupportedFormat:
		return false
	default:
		return true
	}
}
