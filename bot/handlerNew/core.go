package handlerNew

import "log"

func getLogger(scope string) func(s string) {
    return func(s string) {
        log.Printf("Handler HandleNew: [%s] %s", scope, s)
    }
}
