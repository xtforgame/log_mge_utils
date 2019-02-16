package listenert1

import (
	// "errors"
	// "io"
	// "path/filepath"
	"github.com/xtforgame/log_mge_utils/lmu"
)

type ListenerOptionsT1 struct {
	EventReceiverSize int
	BufferSize        int
}

type ListenerT1 struct {
	owner         lmu.Logger
	mode          string
	options       ListenerOptionsT1
	isListening   bool
	isRestoring   bool
	receiveChan   chan *lmu.LoggerEvent
	eventCallback lmu.EventCallback

	logStorerReader lmu.SReader
	logBuffReader   lmu.SReader
	currentPos      int64
	lastEOFPos      int64
}

func CreateListenerT1(logger lmu.Logger, options interface{}) (lmu.Listener, error) {
	listener := &ListenerT1{owner: logger}
	listener.Reinit(options)
	return listener, nil
}

func (listener *ListenerT1) GetOwner() lmu.Logger {
	return listener.owner
}

func (listener *ListenerT1) GetRef() interface{} {
	return listener
}

func (listener *ListenerT1) Reinit(options interface{}) error {
	if options != nil {
		// op, _ := options.(ListenerOptionsT1)
		listener.mode = lmu.ListenerModeCallback
		listener.receiveChan = make(chan *lmu.LoggerEvent, 100)
	}
	listener.closeReaders()
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
		data := &lmu.DataEventPayload{
			IsFromRestoring: true,
			Bytes:           readBuff[:n],
		}
		listener.Dispatch(&lmu.LoggerEvent{
			Name:     lmu.EventOnData,
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

func (listener *ListenerT1) Dispatch(event *lmu.LoggerEvent) {
	// if len(listener.receiveChan) == cap(listener.receiveChan) {
	// }
	if listener.eventCallback != nil {
		listener.eventCallback(event)
	}
}

func (listener *ListenerT1) Receive(event *lmu.LoggerEvent) {
	if !listener.isListening {
		return
	}
	if event.Name == lmu.EventNextIteration {
		listener.Reinit(nil)
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

func (listener *ListenerT1) OnEvent(eventCallback lmu.EventCallback) {
	listener.eventCallback = eventCallback
}

func (listener *ListenerT1) closeReaders() {
	if listener.logStorerReader != nil {
		listener.logStorerReader.Close()
		listener.logStorerReader = nil
	}

	if listener.logBuffReader != nil {
		listener.logBuffReader.Close()
		listener.logBuffReader = nil
	}
}

func (listener *ListenerT1) Close() {
	listener.owner.RemoveListener(listener)

	if listener.receiveChan != nil {
		close(listener.receiveChan)
		listener.receiveChan = nil
	}
	listener.closeReaders()
}
