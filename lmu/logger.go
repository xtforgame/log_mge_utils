package lmu

import (
// "errors"
// "io"
// "path/filepath"
)

type LoggerT1 struct {
	streamSize int64
	logStorer  LogStorer
	logBuffer  LogBuffer
	listeners  []Listener
}

func NewLoggerT1(logStorer LogStorer, logBuffer LogBuffer) (*LoggerT1, error) {
	return &LoggerT1{
		logStorer: logStorer,
		logBuffer: logBuffer,
	}, nil
}

func (logger *LoggerT1) GetPath() string {
	return ""
}

func (logger *LoggerT1) GetIteration() string {
	return ""
}

func (logger *LoggerT1) GetStreamName(iteration string) string {
	return ""
}

func (logger *LoggerT1) GetLogStorer() LogStorer {
	return logger.logStorer
}

func (logger *LoggerT1) GetLogBuffer() LogBuffer {
	return logger.logBuffer
}

func (logger *LoggerT1) Write(p []byte) (int, error) {
	n, err := logger.logStorer.Write(p)
	data := &DataEventPayload{
		Bytes: p,
	}
	for _, listener := range logger.listeners {
		listener.Receive(&LoggerEvent{
			Name:     "data",
			Position: logger.streamSize,
			Data:     data,
		})
	}
	logger.streamSize += int64(n)
	return n, err
}

func (logger *LoggerT1) CreateListener() (Listener, error) {
	listener := &ListenerT1{
		owner:       logger,
		mode:        ListenerModeCallback,
		receiveChan: make(chan *LoggerEvent, 100),
	}
	// err := listener.Reload()
	// if err != nil {
	// 	return nil, err
	// }
	logger.listeners = append(logger.listeners, listener)
	listener.Reinit()
	return listener, nil
}

func (logger *LoggerT1) SwitchToNextIteration(iteration string) (err error) {
	return logger.logStorer.SwitchToNextIteration(iteration)
}

func (logger *LoggerT1) Close() {
	if logger.logStorer != nil {
		logger.logStorer.Close()
		logger.logStorer = nil
	}
}
