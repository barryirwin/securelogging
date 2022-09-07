package logging

import (
	"fmt"

	"go.uber.org/zap"
)

// NewLogger :
//
// Returns a configured zap logger, on the desired level.
// If the loglevel is not debug, info, warn or error, defaults to info.
// Is a wrapper of the other New$LevelLogger() in common
func NewLogger(loglevel string) (*zap.Logger, error) {
	switch loglevel {
	case "debug":
		logger, err := NewDebugLogger()
		if err != nil {
			fmt.Println("NewLogger : Error creating the debug logger:", err.Error())
			return nil, err
		}
		return logger, nil
	case "info":
		logger, err := NewInfoLogger()
		if err != nil {
			fmt.Println("NewLogger : Error creating the info logger:", err.Error())
			return nil, err
		}
		return logger, nil
	case "warn":
		logger, err := NewWarnLogger()
		if err != nil {
			fmt.Println("NewLogger : Error creating the warn logger:", err.Error())
			return nil, err
		}
		return logger, nil
	case "error":
		logger, err := NewErrorLogger()
		if err != nil {
			fmt.Println("NewLogger : Error creating the error logger:", err.Error())
			return nil, err
		}
		return logger, nil
	default:
		logger, err := NewInfoLogger()
		if err != nil {
			fmt.Println("NewLogger : Error creating the default logger:", err.Error())
			return nil, err
		}
		return logger, nil
	}

}
