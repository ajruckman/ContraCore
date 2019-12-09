package db

import (
    "fmt"
    "testing"
)

func TestGetLeaseDetails(t *testing.T) {
    leases := GetLeaseDetails()

    for _, v := range leases {
        fmt.Println(v.IP, v.Hostname, v.Vendor, v.MAC)
    }
}
