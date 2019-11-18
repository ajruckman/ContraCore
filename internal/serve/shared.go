package serve

import (
    "strings"

    "github.com/coredns/coredns/plugin/pkg/log"
)

var (
    clog = log.NewWithPlugin("contradomain")
)

func rt(in string) string {
    return strings.TrimSuffix(in, ".")
}
