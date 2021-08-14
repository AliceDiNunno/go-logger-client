package GoLoggerClient

import (
	"fmt"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

type GoLoggerHook struct {
	config ClientConfiguration
}

func (l GoLoggerHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

func getHostname() string {
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "Unknown"
	}
	return hostname
}

func getNestedTraces(errorTrace *e.Error) []*Traceback {
	var nestedTrace []*Traceback

	for {
		newTrace := &Traceback{
			Message:   errorTrace.Err.Error(),
			Traceback: nil,
		}

		for _, trace := range errorTrace.Stack {
			newTrace.Traceback = append(newTrace.Traceback, e.Frame{
				Filename: trace.Filename,
				Method:   trace.Method,
				Line:     trace.Line,
			})
		}

		nestedTrace = append(nestedTrace, newTrace)

		errorTrace = errorTrace.Child
		if errorTrace == nil {
			break
		}
	}

	return nestedTrace
}

func (l GoLoggerHook) Fire(entry *logrus.Entry) error {
	var fields map[string]interface{}

	fields = entry.Data
	errorData := fields["err"]
	errorTrace := *(errorData.(**e.Error))
	fingerprint := errorTrace.Fingerprint()
	resultingError := errorTrace.Err
	nestedTrace := getNestedTraces(errorTrace)

	var resultingTrace *Traceback = nil

	if nestedTrace != nil && len(nestedTrace) > 0 {
		resultingTrace = nestedTrace[0]
	}

	data := ItemCreationRequest{
		ProjectKey: l.config.Key,
		Identification: LogIdentification{
			Client: LogClientIdentification{
				UserID:    nil,
				IPAddress: "",
			},
			Deployment: LogDeploymentIdentification{
				Platform:    "unix",
				Source:      "server",
				Hostname:    getHostname(),
				Environment: l.config.Environment,
				Version:     l.config.Version,
			},
		},
		Data: LogData{
			Timestamp:        time.Now(),
			GroupingID:       fingerprint,
			Fingerprint:      fingerprint,
			Level:            entry.Level.String(),
			Trace:            resultingTrace,
			NestedTrace:      nestedTrace,
			Message:          resultingError.Error(),
			StatusCode:       -1,
			AdditionalFields: nil,
		},
	}

	delete(fields, "err")

	if value, hasValue := entry.Data["code"]; hasValue {
		status, err := strconv.Atoi(fmt.Sprintf("%v", value))
		if err == nil {
			data.Data.StatusCode = status
		}
		delete(fields, "code")
	}

	if value, hasValue := entry.Data["user"]; hasValue {
		userId, err := uuid.Parse(fmt.Sprintf("%v", value))
		if err == nil {
			data.Identification.Client.UserID = &userId
		}
		delete(fields, "user")
	}

	if value, hasValue := entry.Data["ip"]; hasValue {
		ip := fmt.Sprintf("%v", value)

		data.Identification.Client.IPAddress = ip

		delete(fields, "ip")
	}

	data.Data.AdditionalFields = fields

	err := l.sendDataToServer(data)

	return err
}

func SetupHook(config ClientConfiguration) {
	hook := &GoLoggerHook{
		config: config,
	}
	logrus.AddHook(hook)
}
