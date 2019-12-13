package schema

import (
    "net"
    "time"
)

type LeaseDetails struct {
    Time     time.Time
    Op       string
    MAC      string
    IP       net.IP
    Hostname string
    Vendor   string
}
