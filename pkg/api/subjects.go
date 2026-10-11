package api

const (
	RegistrationSubject = "magnetar.registration"
	NodeUpdateSubject   = "magnetar.node.update"
	// VmDesiredGetSubject is how a host's Star polls the desired VM set
	// (standing declaration stored on the node record). VmStatusSetSubject
	// is how the host reports observed VM status back.
	VmDesiredGetSubject = "magnetar.vms.desired.get"
	VmStatusSetSubject  = "magnetar.vms.status.set"
	OsDesiredGetSubject = "magnetar.os.desired.get"
)
