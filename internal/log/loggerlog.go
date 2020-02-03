package log

import (
    "go.uber.org/atomic"
)

var (
    AnsweredTotDuration atomic.Duration
    AnsweredTotCount    atomic.Uint32
    PassedTotDuration   atomic.Duration
    PassedTotCount      atomic.Uint32

    LogDurations = true
)
