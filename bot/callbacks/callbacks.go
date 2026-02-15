package callbacks

import (
    "useful.team/bloodpressure/m/bot/core"
)

var (
    Callbacks = make(map[string]core.Handler)
)

func GetHandler(data string) (core.Handler, bool) {
    if handler, found := Callbacks[data]; found {
        return handler, true
    }

    return nil, false
}

func RegisterCallback(id string, f core.Handler) {
    Callbacks[id] = f
}
