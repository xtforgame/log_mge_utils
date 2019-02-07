package lmu

import (
	"errors"
)

var ITER_EXPIRED = errors.New("iteration expired")
var EventNextIteration = "next-iteration"
var ListenerModeCallback = "listener-mode-callback"
var ListenerModeRead = "listener-mode-read"
