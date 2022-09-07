package logging

import (
	"encoding/json"

	"go.uber.org/zap"
)

// NewInfoLogger :
//
// Returns a configured zap logger, on Info level
func NewInfoLogger() (*zap.Logger, error) {
	rawJSONConfig := []byte(`{
    "level": "info",
    "encoding": "console",
    "outputPaths": ["stdout"],
    "errorOutputPaths": ["stdout"],
    "encoderConfig": {
        "messageKey": "message",
        "levelKey": "level",
        "nameKey": "logger",
        "timeKey": "time",
        "callerKey": "logger",
        "stacktraceKey": "stacktrace",
        "callstackKey": "callstack",
        "errorKey": "error",
        "timeEncoder": "iso8601",
        "fileKey": "file",
        "levelEncoder": "capitalColor",
        "durationEncoder": "second",
        "callerEncoder": "short",
        "nameEncoder": "full",
        "sampling": {
            "initial": "3",
            "thereafter": "10"
            }
        }
    }`)
	config := zap.Config{}
	if err := json.Unmarshal(rawJSONConfig, &config); err != nil {
		return nil, err
	}
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}
