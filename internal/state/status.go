package state

import (
    "go.uber.org/atomic"
)

var PostgresOnline atomic.Bool

var ClickHouseOnline atomic.Bool

