package handlerStart

import "log"

func getLogger(scope string) func(s string) {
    return func(s string) {
        log.Printf("Handler Start: [%s] %s", scope, s)
    }
}
