package system

import (
    "go.uber.org/atomic"
)

var (
    ContraDBOnline  atomic.Bool
    ContraLogOnline atomic.Bool
)
