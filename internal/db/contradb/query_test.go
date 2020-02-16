package contradb

import (
    "fmt"
    "testing"

    "github.com/stretchr/testify/assert"

    "github.com/ajruckman/ContraCore/internal/config"
)

func init() {
    config.ContraDBURL = "postgres://contradbmgr:contradbmgr@127.0.0.1/contradb?timezone=UTC"
    Setup()
}

func TestGetLeaseDetails(t *testing.T) {
    leases, err := GetLeaseDetails()
    assert.Equal(t, err, nil)

    for _, v := range leases {
        fmt.Println(v.IP, v.Hostname, v.Vendor, v.MAC)
    }
}

func TestGetWhitelistRules(t *testing.T) {
    _, err := Exec(`truncate table whitelist;`)
    assert.Equal(t, err, nil)

    whitelist, err := GetWhitelistRules()
    assert.Equal(t, err, nil)
    assert.Empty(t, whitelist)

    // All exclusions null
    _, err = Exec(`insert into whitelist (pattern) values ('(?:^|.*\.)google.com')`)
    assert.NotEqual(t, err, nil)

    _, err = Exec(`insert into whitelist (pattern, ips) values ('(?:^|.*\.)google.com', '{ 1.1.1.1, 1.0.0.1 }')`)
    assert.Equal(t, err, nil)

    _, err = Exec(`insert into whitelist (pattern, subnets) values ('(?:^|.*\.)google.com', '{ 10.1.7.0/24, 10.5.0.0/16 }')`)
    assert.Equal(t, err, nil)

    _, err = Exec(`insert into whitelist (pattern, macs) values ('(?:^|.*\.)google.com', '{ aa:ff:00:11:22:dd, bb:cc:00:33:99:cc }')`)
    assert.Equal(t, err, nil)

    _, err = Exec(`insert into whitelist (pattern, vendors) values ('(?:^|.*\.)google.com', '{ ".*Apple.*", "^Xiaomi.*" }')`)
    assert.Equal(t, err, nil)

    _, err = Exec(`insert into whitelist (pattern, hostnames) values ('(?:^|.*\.)google.com', '{ "a7zkk3lu", "378xkjql.2jd.jik.c" }')`)
    assert.Equal(t, err, nil)

    _, err = Exec(`insert into whitelist (pattern, subnets) values ('(?:^|.*\.)badsite.biz', '{ 0.0.0.0/0 }')`)
    assert.Equal(t, err, nil)

    whitelist, err = GetWhitelistRules()
    assert.Equal(t, err, nil)
    assert.NotEmpty(t, whitelist)

    assert.Nil(t, whitelist[0].Subnets)
    assert.Nil(t, whitelist[1].MACs)
    assert.Nil(t, whitelist[2].Vendors)
    assert.Nil(t, whitelist[3].Hostnames)
    assert.Nil(t, whitelist[4].IPs)

    for _, v := range whitelist {
        assert.Nil(t, v.Expires)
        fmt.Println(v)
    }
}
