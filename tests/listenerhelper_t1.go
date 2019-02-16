package tests

import (
	"testing"
	// "errors"
	// "bufio"
	"github.com/xtforgame/log_mge_utils/lmu"
	// "github.com/xtforgame/log_mge_utils/logbuffers"
	// "github.com/xtforgame/log_mge_utils/logstorers/localfs"
	// "os"
)

type ListenerHelperT1 struct {
	t        *testing.T
	Listener lmu.Listener

	Name                string
	RestoreEventCounter int
	WriteEventCounter   int
}

func CreateListenerHelperT1(t *testing.T, name string, listener lmu.Listener) *ListenerHelperT1 {
	lh := &ListenerHelperT1{
		t:        t,
		Name:     name,
		Listener: listener,
	}
	listener.OnEvent(func(event *lmu.LoggerEvent) {
		data, ok := event.Data.(*lmu.DataEventPayload)
		if event.Name == lmu.EventOnData && ok {
			if data.IsFromRestoring {
				lh.RestoreEventCounter++
			} else {
				lh.WriteEventCounter++
			}
			t.Log(lh.Name+": data.IsFromRestoring", data.IsFromRestoring)
			t.Log(lh.Name+": event.Position", event.Position)
			t.Log("")
		}
	})
	return lh
}

func (lh *ListenerHelperT1) AssertRestoreEventCounter(counter int) {
	if lh.RestoreEventCounter != counter {
		lh.t.Log("CurrentPos of "+lh.Name+": ", lh.Listener.GetCurrentPos())
		lh.t.Fatal("expect RestoreEventCounter of "+lh.Name+": ", counter, ", actual:", lh.RestoreEventCounter)
	}
}

func (lh *ListenerHelperT1) AssertWriteEventCounter(counter int) {
	if lh.WriteEventCounter != counter {
		lh.t.Log("CurrentPos of "+lh.Name+": ", lh.Listener.GetCurrentPos())
		lh.t.Fatal("expect WriteEventCounter of "+lh.Name+": ", counter, ", actual:", lh.WriteEventCounter)
	}
}
