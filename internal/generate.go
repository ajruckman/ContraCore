package internal

//go:generate go run ./provision/generate.go
//go:generate sqlboiler --wipe -o contradb -p contradb psql
