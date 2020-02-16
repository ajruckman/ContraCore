package system

import (
    "sync"
)

var (
    HasInitializedLock sync.Mutex
    HasInitialized     bool
)
