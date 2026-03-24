package handlerLog

import (
    "log"
    "strings"
)

func getLogger(scope string) func(s ...string) {
    return func(s ...string) {
        log.Printf("Handler LogShort: [%s] %s", scope, strings.Join(s, " "))
    }
}
