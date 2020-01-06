package serve

import (
    "sync"

    "github.com/ajruckman/ContraCore/internal/db"

    . "github.com/ajruckman/xlib"
)

var (
    ouiMACPrefixToVendor sync.Map
)

func cacheOUI() {
    oui, err := db.GetOUI()
    Err(err)

    for _, v := range oui {
        ouiMACPrefixToVendor.Store(v.MAC, v.Vendor)
    }
}
