package serve

import (
    "strings"
)

func rt(in string) string {
    return strings.TrimSuffix(in, ".")
}
