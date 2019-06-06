package logwatcher

import (
	// "fmt"
	// "flag"
	"net/http"
	// "strings"
	"github.com/go-chi/chi"
	// "github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/xtforgame/log_mge_utils/lmu"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WsReaderEvent struct {
	Bytes    []byte
	Error    error
	Finished bool
}

func basicHandler(logID string, conn *websocket.Conn, logger lmu.Logger) {
	var listener lmu.Listener
	defer func() {
		if listener != nil {
			listener.Close()
		}
	}()

	readChan := make(chan WsReaderEvent)

	go func() {
		for {
			_, p, err := conn.ReadMessage()
			readChan <- WsReaderEvent{
				Bytes: p,
				Error: err,
			}
			if err != nil {
				// fmt.Println("ws:", err)
				return
			}
		}
	}()

	for {
		// messageType, p, err := conn.ReadMessage()
		readEvent := <-readChan
		if readEvent.Error != nil || readEvent.Finished {
			// fmt.Println("ws:", err)
			// fmt.Println("readEvent.Finished", readEvent.Finished)
			return
		}
		p := readEvent.Bytes

		if len(p) > 0 {
			if listener == nil {
				l, _ := logger.CreateListener(nil)
				listener = l
				l.OnEvent(func(event *lmu.LoggerEvent) {
					if event.Name == lmu.EventOnData {
						data, ok := event.Data.(*lmu.DataEventPayload)
						if ok {
							bytes := append([]byte{EventOnDataCode, 0}, data.Bytes...)
							// fmt.Println("data :", data)
							// if err := conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
							if err := conn.WriteMessage(websocket.BinaryMessage, bytes); err != nil {
								// log.Println(err)
								return
							}
						}
					} else if event.Name == lmu.EventNextIteration {
						bytes := []byte{EventNextIterationCode}
						if err := conn.WriteMessage(websocket.BinaryMessage, bytes); err != nil {
							// log.Println(err)
							return
						}
					} else if event.Name == lmu.EventLogRemoved {
						bytes := []byte{EventLogRemovedCode}
						if err := conn.WriteMessage(websocket.BinaryMessage, bytes); err != nil {
							// log.Println(err)
						}
						// fmt.Println("lmu.EventLogRemoved")
						go func() {
							readChan <- WsReaderEvent{
								Finished: true,
							}
						}()
					}
				})
				l.Restore()
				l.Listen()
			} else if p[0] == 1 {
				logger.Write(p[1:])
			} else if p[0] == 2 {
				logger.SwitchToNextIteration("")
			} else if p[0] == 3 {
				LoggerHeplerInst.RemoveAndCloseLogger(logID)
			}
		}
		// fmt.Println("p :", p)
		// conn.WriteMessage(websocket.TextMessage, []byte(""))
		// conn.WriteMessage(websocket.BinaryMessage, []byte(""))
		// if err := conn.WriteMessage(messageType, p); err != nil {
		// 	// log.Println(err)
		// 	return
		// }
	}
}

func LoggerWebsocket(w http.ResponseWriter, r *http.Request) {
	// l := log.WithField("remoteaddr", r.RemoteAddr)
	logID := chi.URLParam(r, "logID")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// l.WithError(err).Error("Unable to upgrade connection")
		return
	}

	defer func() {
		conn.Close()
	}()

	logger := LoggerHeplerInst.CreateOrGetLogger(logID)
	if logger == nil {
		bytes := append([]byte{WrongPathCode}, []byte("Wrong path name: "+logID)...)
		conn.WriteMessage(websocket.BinaryMessage, bytes)
		return
	}
	basicHandler(logID, conn, logger)
}

func ListenerWebsocket(w http.ResponseWriter, r *http.Request) {
	// l := log.WithField("remoteaddr", r.RemoteAddr)
	logID := chi.URLParam(r, "logID")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// l.WithError(err).Error("Unable to upgrade connection")
		return
	}

	defer func() {
		conn.Close()
	}()

	logger := LoggerHeplerInst.CreateOrGetLogger(logID)
	if logger == nil {
		bytes := append([]byte{WrongPathCode}, []byte("Wrong path name: "+logID)...)
		conn.WriteMessage(websocket.BinaryMessage, bytes)
		return
	}
	basicHandler(logID, conn, logger)
}
