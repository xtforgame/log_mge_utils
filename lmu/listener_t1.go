package lmu

import (
// "errors"
// "io"
// "path/filepath"
)

type ListenerOptionsT1 struct {
	EventReceiverSize int
	BufferSize        int
}

type ListenerT1 struct {
	owner         *LoggerT1
	mode          string
	options       ListenerOptionsT1
	isListening   bool
	isRestoring   bool
	receiveChan   chan *LoggerEvent
	eventCallback EventCallback

	logStorerReader SReader
	logBuffReader   SReader
	currentPos      int64
	lastEOFPos      int64
}

func (listener *ListenerT1) GetOwner() Logger {
	return listener.owner
}

func (listener *ListenerT1) Reinit() error {
	newStorerReader, err := listener.owner.GetLogStorer().CreateReader()
	if err != nil {
		return err
	}
	listener.logStorerReader = newStorerReader

	newBuffReader, err := listener.owner.GetLogBuffer().CreateReader()
	if err != nil {
		return err
	}
	listener.logBuffReader = newBuffReader
	return nil
}

func (listener *ListenerT1) GetIteration() string {
	return listener.GetOwner().GetIteration()
}

func (listener *ListenerT1) GetCurrentPos() int64 {
	// pos, _ := listener.file.Seek(0, io.SeekCurrent)
	return listener.currentPos
}

func (listener *ListenerT1) GetLastEOFPos() int64 {
	return listener.lastEOFPos
}

func (listener *ListenerT1) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (listener *ListenerT1) Seek(off int64, whence int) (int64, error) {
	return 0, nil
}

func (listener *ListenerT1) ReloadAndRead(p []byte) (int, error) {
	return 0, nil
}

func (listener *ListenerT1) Restore() error {
	readBuff := make([]byte, 200)
	listener.isRestoring = true
	for true {
		n, err := listener.logStorerReader.Read(readBuff)
		if err != nil {
			listener.isRestoring = false
			return err
		}
		listener.currentPos += int64(n)
		data := &DataEventPayload{
			IsFromRestoring: true,
			Bytes:           readBuff[:n],
		}
		listener.Dispatch(&LoggerEvent{
			Name:     "data",
			Position: listener.currentPos,
			Data:     data,
		})
	}
	return nil
}

func (listener *ListenerT1) Listen() {
	listener.isListening = true
}

func (listener *ListenerT1) Unlisten() {
	listener.isListening = false
}

func (listener *ListenerT1) Dispatch(event *LoggerEvent) {
	if len(listener.receiveChan) == cap(listener.receiveChan) {

	}
	if listener.eventCallback != nil {
		listener.eventCallback(event)
	}
}

func (listener *ListenerT1) Receive(event *LoggerEvent) {
	if !listener.isListening {
		return
	}
	if event.Name == EventNextIteration {
		listener.Reinit()
		listener.Dispatch(event)
		// listener.isRestoring = true
		return
	}
	if listener.isRestoring {
		return
	}
	if listener.currentPos < event.Position {
		listener.Restore()
		return
	}
	listener.currentPos += event.Length
	listener.Dispatch(event)
}

func (listener *ListenerT1) OnEvent(eventCallback EventCallback) {
	listener.eventCallback = eventCallback
}

func (listener *ListenerT1) Close() {
	if listener.logStorerReader != nil {
		listener.logStorerReader.Close()
		listener.logStorerReader = nil
	}

	if listener.logBuffReader != nil {
		listener.logBuffReader.Close()
		listener.logBuffReader = nil
	}
	if listener.receiveChan != nil {
		close(listener.receiveChan)
		listener.receiveChan = nil
	}
}
