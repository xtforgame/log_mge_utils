package logbuffers

import (
	"github.com/xtforgame/log_mge_utils/lmu"
)

type SimpleBuffer struct {
	offset int64
}

func NewSimpleBuffer() (*SimpleBuffer, error) {
	return &SimpleBuffer{
		offset: 0,
	}, nil
}

func (lb *SimpleBuffer) GetPath() string {
	return ""
}

func (lb *SimpleBuffer) GetIteration() string {
	return ""
}

func (lb *SimpleBuffer) GetStreamName(iteration string) string {
	return ""
}

func (lb *SimpleBuffer) SwitchToNextIteration(iteration string) error {
	return nil
}

func (lb *SimpleBuffer) Write(p []byte) (int, error) {
	return 0, nil
}

func (lb *SimpleBuffer) CreateReader() (lmu.SReader, error) {
	return nil, nil
}

func (lb *SimpleBuffer) GetOffset() int64 {
	return 0
}

func (lb *SimpleBuffer) Forget(offset int64) int64 {
	return 0
}

func (lb *SimpleBuffer) RemoveBuffer() {

}

func (lb *SimpleBuffer) Close() {

}
