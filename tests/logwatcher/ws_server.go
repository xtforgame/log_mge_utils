package logwatcher

import (
	"fmt"
	// "flag"
	"net/http"
	// "strings"
	// "github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/xtforgame/log_mge_utils/lmu"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func TestHandleWebsocket(w http.ResponseWriter, r *http.Request) {
	// l := log.WithField("remoteaddr", r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// l.WithError(err).Error("Unable to upgrade connection")
		return
	}

	defer func() {
		conn.Close()
	}()

	if LoggerHeplerInst == nil {
		LoggerHeplerInst = CreateLoggerHepler()
	}

	var listener lmu.Listener
	l, _ := LoggerHeplerInst.Logger.CreateListener(nil)
	listener = l
	l.OnEvent(func(event *lmu.LoggerEvent) {
		if event.Name == lmu.EventOnData {
			data, ok := event.Data.(*lmu.DataEventPayload)
			if ok {
				bytes := append([]byte{1, 0}, data.Bytes...)
				// fmt.Println("data :", data)
				// if err := conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
				if err := conn.WriteMessage(websocket.BinaryMessage, bytes); err != nil {
					// log.Println(err)
					return
				}
			}
		} else if event.Name == lmu.EventNextIteration {
			bytes := []byte{2, 0}
			if err := conn.WriteMessage(websocket.BinaryMessage, bytes); err != nil {
				// log.Println(err)
				return
			}
		}
	})
	l.Restore()
	l.Listen()
	defer func() {
		listener.Close()
	}()

	for {
		// messageType, p, err := conn.ReadMessage()
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("ws:", err)
			return
		}
		if len(p) > 0 && p[0] == 2 {
			LoggerHeplerInst.Logger.SwitchToNextIteration("")
		} else {
			LoggerHeplerInst.Logger.Write(p)
		}
		fmt.Println("p :", p)
		// conn.WriteMessage(websocket.TextMessage, []byte(""))
		// conn.WriteMessage(websocket.BinaryMessage, []byte(""))
		// if err := conn.WriteMessage(messageType, p); err != nil {
		// 	// log.Println(err)
		// 	return
		// }
	}
}
