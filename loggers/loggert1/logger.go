package loggert1

import (
	// "errors"
	// "io"
	// "path/filepath"
	"github.com/xtforgame/log_mge_utils/lmu"
	"sync"
)

type LoggerT1 struct {
	streamSize     int64
	createListener lmu.ListenerCreator
	logStorer      lmu.LogStorer
	logBuffer      lmu.LogBuffer
	listeners      map[interface{}]lmu.Listener
	listenersMu    sync.Mutex
}

func NewLoggerT1(logStorer lmu.LogStorer, logBuffer lmu.LogBuffer, createListener lmu.ListenerCreator) (*LoggerT1, error) {
	return &LoggerT1{
		logStorer:      logStorer,
		logBuffer:      logBuffer,
		createListener: createListener,
		listeners:      map[interface{}]lmu.Listener{},
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

func (logger *LoggerT1) GetLogStorer() lmu.LogStorer {
	return logger.logStorer
}

func (logger *LoggerT1) GetLogBuffer() lmu.LogBuffer {
	return logger.logBuffer
}

func (logger *LoggerT1) GetStreamSize() int64 {
	return logger.streamSize
}

func (logger *LoggerT1) dispatch(event *lmu.LoggerEvent) {
	logger.listenersMu.Lock()
	for _, listener := range logger.listeners {
		listener.Receive(event)
	}
	logger.listenersMu.Unlock()
}

func (logger *LoggerT1) Write(p []byte) (int, error) {
	n, err := logger.logStorer.Write(p)
	data := &lmu.DataEventPayload{
		Bytes: p,
	}
	logger.dispatch(&lmu.LoggerEvent{
		Name:     lmu.EventOnData,
		Position: logger.streamSize,
		Length:   int64(n),
		Data:     data,
	})
	logger.streamSize += int64(n)
	return n, err
}

func (logger *LoggerT1) CreateListener(options interface{}) (lmu.Listener, error) {
	logger.listenersMu.Lock()
	listener, err := logger.createListener(logger, options)
	if err == nil {
		logger.listeners[listener.GetRef()] = listener
	}
	logger.listenersMu.Unlock()
	return listener, err
}

func (logger *LoggerT1) RemoveListener(listener lmu.Listener) {
	logger.listenersMu.Lock()
	delete(logger.listeners, listener.GetRef())
	logger.listenersMu.Unlock()
}

func (logger *LoggerT1) SwitchToNextIteration(iteration string) error {
	err := logger.logStorer.SwitchToNextIteration(iteration)
	if err == nil {
		logger.dispatch(&lmu.LoggerEvent{
			Name:     lmu.EventNextIteration,
			Position: logger.streamSize,
			Length:   0,
			Data:     nil,
		})
	}
	return err
}

func (logger *LoggerT1) Close() {
	if logger.logStorer != nil {
		logger.logStorer.Close()
		logger.logStorer = nil
	}
}
