// +build ignore

package main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "os"
)

func main() {
    fmt.Println("Generating run.go")

    err := os.Chdir("./provision")
    if err != nil {
        panic(err)
    }

    if _, err = os.Stat("./run.go"); err == nil {
        err = os.Remove("./run.go")
        if err != nil {
            panic(err)
        }
    }

    file, err := os.Create("run.go")
    if err != nil {
        panic(err)
    }

    w := bufio.NewWriter(file)

    ddlsql, err := ioutil.ReadFile("./sql/10_ddl.sql")
    if err != nil {
        panic(err)
    }

    w.WriteString("package provision\n\n")
    w.WriteString("// Code generated by generate.go; DO NOT EDIT.\n\n")
    w.WriteString("import `context`\n")
    w.WriteString("import `github.com/ajruckman/ContraCore/internal/db`\n")
    w.WriteString("import `github.com/coredns/coredns/plugin/pkg/log`\n\n")
    w.WriteString("func Setup() {\n")
    w.WriteString("    clog := log.NewWithPlugin(`contradomain`)\n")
    w.WriteString("    if !db.PostgresOnline.Load() {\n")
    w.WriteString("        clog.Warning(`not connected to PostgreSQL; not provisioning database`)\n")
    w.WriteString("        return\n")
    w.WriteString("    }\n")
    w.WriteString("    clog.Info(`provisioning database`)\n")
    w.WriteString("    _, err := db.PDB.Exec(context.Background(), `\n")
    w.WriteString(string(ddlsql))
    w.WriteString("    `)\n")
    w.WriteString("    if err != nil { panic(err) }\n")
    w.WriteString("}")

    w.Flush()
    file.Close()
}
