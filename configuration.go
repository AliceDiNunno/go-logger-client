package GoLoggerClient

type ClientConfiguration struct {
	Url         string
	ProjectId   string
	Port        int
	Key         string
	Environment string
	Version     string

	RemoveFieldsFromDebugOutput bool
}

type ClientTransporter interface {
	Send(data ItemCreationRequest) error
}
