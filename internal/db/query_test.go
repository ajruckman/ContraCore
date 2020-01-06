package db

import (
    "fmt"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestGetLeaseDetails(t *testing.T) {
    leases, err := GetLeaseDetails()
    assert.Equal(t, err, nil)

    for _, v := range leases {
        fmt.Println(v.IP, v.Hostname, v.Vendor, v.MAC)
    }
}
