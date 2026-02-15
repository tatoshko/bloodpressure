package pgsql

import (
    "fmt"
    "strings"
)

func getLogger(scope string) func(s ...string) {
    return func(s ...string) {
        fmt.Printf("PG [%s] %s", scope, strings.Join(s, ""))
    }
}
