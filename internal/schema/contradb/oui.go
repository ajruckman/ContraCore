package contradb

type OUI struct {
    MAC    string `db:"mac"`
    Vendor string `db:"vendor"`
}
