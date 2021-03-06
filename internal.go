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

	eerr := t.Receiver.PushNewLogEntry(id, &creationRequest)

	if eerr != nil {
		return eerr.Err
	}

	return nil
}

func NewInternalTransporter(receiver InternalReceiver, config ClientConfiguration) *InternalTransporter {
	return &InternalTransporter{
		Config:   config,
		Receiver: receiver,
	}
}
