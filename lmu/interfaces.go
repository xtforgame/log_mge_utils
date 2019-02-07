package lmu

type Stream interface {
	GetIteration() (iteration string)
	Close()
}

type Writer interface {
	Stream
	GetPath() (path string)
	GetStreamName(iteration string) (name string)
	SwitchToNextIteration(iteration string) (err error)
	Write(p []byte) (n int, err error)
}

type Readable interface {
	Stream
	GetOwner() Writer
	GetLastEOFPos() (n int64)
	GetCurrentPos() (n int64)
}

type EventCallback func(event *LoggerEvent)

type Listener interface {
	Readable
	StartRestore() error
	Receive(event *LoggerEvent)
	OnEvent(eventCallback EventCallback)
}

type SReader interface {
	Readable
	Read(p []byte) (n int, err error)
	ReloadAndRead(p []byte) (n int, err error)
	Seek(off int64, whence int) (ret int64, err error)
}

type ReadableWriter interface {
	Writer
	CreateReader() (SReader, error)
}

type LogStorer interface {
	ReadableWriter
	GetStoreType() (_type string)
}

type LogBuffer interface {
	ReadableWriter
	GetOffset() (offset int64)
	Forget(offset int64) (newOffset int64)
}
