package GoLoggerClient

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/google/uuid"
)

type InternalReceiver interface {
	PushNewLogEntry(id uuid.UUID, request *ItemCreationRequest) *e.Error
}

type InternalTransporter struct {
	config   ClientConfiguration
	receiver InternalReceiver
}

func (t InternalTransporter) Send(creationRequest ItemCreationRequest) error {
	id, err := uuid.Parse(t.config.ProjectId)

	if err != nil {
		return err
	}

	return t.receiver.PushNewLogEntry(id, &creationRequest).Err
}

func NewInternalTransporter(receiver InternalReceiver, config ClientConfiguration) *InternalTransporter {
	return &InternalTransporter{
		config:   config,
		receiver: receiver,
	}
}
