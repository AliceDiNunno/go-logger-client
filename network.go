package go_logger_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (l GoLoggerHook) sendDataToServer(creationRequest ItemCreationRequest) error {
	bodyJson, err := json.Marshal(creationRequest)

	if err != nil {
		return err
	}

	bodyBuf := bytes.NewBuffer(bodyJson)

	request, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%d/%s/items", l.config.Url, l.config.Port, l.config.ProjectId), bodyBuf)

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
