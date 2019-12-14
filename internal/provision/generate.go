// +build ignore

package main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "os"
)

func main() {
    fmt.Println("Generating provision.go")

    if _, err := os.Stat("./provision.go"); err == nil {
        err = os.Remove("./provision.go")
        if err != nil {
            panic(err)
        }
    }

    file, err := os.Create("provision.go")
    if err != nil {
        panic(err)
    }

    w := bufio.NewWriter(file)

    ddlsql, err := ioutil.ReadFile("./sql/10_ddl.sql")
    if err != nil {
        panic(err)
    }

    w.WriteString("package db\n\n")
    w.WriteString("import `fmt`\n\n")
    w.WriteString("func init() {\n")
    w.WriteString("    fmt.Println(`Provisioning database`)\n")
    w.WriteString("    _, err := PDB.Exec(`\n")
    w.WriteString(string(ddlsql))
    w.WriteString("    `)\n")
    w.WriteString("    if err != nil { panic(err) }\n")
    w.WriteString("}")

    w.Flush()
    file.Close()
}
