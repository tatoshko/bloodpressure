package handlerLog

import "log"

func getLogger(scope string) func(s string) {
    return func(s string) {
        log.Printf("Handler LogShort: [%s] %s", scope, s)
    }
}
