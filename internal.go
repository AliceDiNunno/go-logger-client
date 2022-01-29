package GoLoggerClient

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/google/uuid"
)

type InternalReceiver interface {
	PushNewLogEntry(id uuid.UUID, request *ItemCreationRequest) *e.Error
}

type InternalTransporter struct {
	Config   ClientConfiguration
	Receiver InternalReceiver
}

func (t InternalTransporter) Send(creationRequest ItemCreationRequest) error {
	id, err := uuid.Parse(t.Config.ProjectId)

	if err != nil {
		return err
	}

	return t.Receiver.PushNewLogEntry(id, &creationRequest).Err
}

func NewInternalTransporter(receiver InternalReceiver, config ClientConfiguration) *InternalTransporter {
	return &InternalTransporter{
		Config:   config,
		Receiver: receiver,
	}
}
