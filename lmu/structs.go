package lmu

type LoggerEvent struct {
	Name     string
	Position int64
	Length   int64
	Data     interface{}
}

type DataEventPayload struct {
	IsFromRestoring bool
	Bytes           []byte
}
