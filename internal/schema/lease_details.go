package schema

import (
    "net"
    "time"
)

type LeaseDetails struct {
    ID       uint64
    Time     time.Time
    Op       string
    MAC      string
    IP       net.IP
    Hostname string
    Vendor   string
}
