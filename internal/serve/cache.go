package serve

import (
    "strings"

    "github.com/ajruckman/ContraCore/internal/db"
    "github.com/ajruckman/ContraCore/internal/schema"
)

var dhcpCacheMap = map[string][]schema.LeaseDetails{}

func dhcpCache() {
    leases := db.GetLeaseDetails()

    for _, lease := range leases {
        hostname := strings.ToLower(lease.Hostname)
        if _, exists := dhcpCacheMap[hostname]; !exists {
            dhcpCacheMap[hostname] = []schema.LeaseDetails{}
        }
        dhcpCacheMap[hostname] = append(dhcpCacheMap[hostname], lease)
    }
}
