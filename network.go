package GoLoggerClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type NetworkTransporter struct {
	config ClientConfiguration
}

func (t NetworkTransporter) Send(creationRequest ItemCreationRequest) error {
	bodyJson, err := json.Marshal(creationRequest)

	if err != nil {
		return err
	}

	bodyBuf := bytes.NewBuffer(bodyJson)

	request, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%d/%s/items", t.config.Url, t.config.Port, t.config.ProjectId), bodyBuf)

	if err != nil {
		request = nil
		return err
	}

	request.Header.Add("content-type", "application/json")
	request.Header.Add("Accept", "*/*")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	bodyBuf = nil
	bodyJson = nil

	if err != nil {
		request = nil
		return err
	}
	_, err = client.Do(request)
	request = nil

	return err
}

func NewNetworkTransporter(config ClientConfiguration) *NetworkTransporter {
	return &NetworkTransporter{
		config: config,
	}
}
