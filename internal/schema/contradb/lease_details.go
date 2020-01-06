package contradb

import (
    "net"
    "time"
)

type LeaseDetails struct {
    Time     time.Time
    Op       string
    IP       net.IP
    MAC      net.HardwareAddr
    Hostname *string
    Vendor   *string
}
