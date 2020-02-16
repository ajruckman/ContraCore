package dbschema

type OUI struct {
	MAC    string `db:"mac"`
	Vendor string `db:"vendor"`
}
