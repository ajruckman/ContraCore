package serve

import (
    "strings"

    . "github.com/ajruckman/xlib"

    "github.com/ajruckman/ContraCore/internal/db"
    "github.com/ajruckman/ContraCore/internal/schema"
)

var DHCPCache = map[string][]schema.LeaseDetails{}

func cacheDHCP() {
    leases, err := db.GetLeaseDetails()
    Err(err)

    for _, lease := range leases {
        hostname := strings.ToLower(lease.Hostname)
        if _, exists := DHCPCache[hostname]; !exists {
            DHCPCache[hostname] = []schema.LeaseDetails{}
        }
        DHCPCache[hostname] = append(DHCPCache[hostname], lease)
    }
}
