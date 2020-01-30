package internal

//go:generate go run ./provision/generate.go
//xgo:generate sqlboiler --wipe -o contradb -p contradb psql
