package main

import (
  "database/sql"
  "flag"
  "fmt"
  "os"
  "strings"
  _ "github.com/lib/pq"
)

type Args struct {
  file *string
  table *string
  column *string
  dbconn *string
}

const pgProtocol = "postgresql://"

func main() {
  fmt.Println("Import data from a file to a single column table.")

  var args Args
  args.file = flag.String("file", "", "source file")
  args.table = flag.String("table", "", "target table")
  args.column = flag.String("column", "", "target column")
  args.dbconn = flag.String("dbconn", "root:secret@localhost:5432/namesdb?sslmode=disable", "postgres connect string")

  flag.Parse()
  fmt.Println("\tfile", *args.file)
  fmt.Println("\ttable", *args.table)
  fmt.Println("\tcolumn", *args.column)
  fmt.Println("\tdbconn", *args.dbconn)

  if *args.file == "" || *args.table == "" || *args.column == "" || *args.dbconn == "" {
    fmt.Fprintf(os.Stderr, "missing arguments, use --help flag for more info.\n")
    os.Exit(1)
  }

  if !strings.HasPrefix(*args.dbconn, pgProtocol) {
    *args.dbconn = pgProtocol + *args.dbconn
    fmt.Printf("Warning: added PostgreSQL protocol to dbconn: %v.\n", *args.dbconn)
  }

  db, err := sql.Open("postgres", *args.dbconn)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error: failed to open the DB connection.\n")
    os.Exit(1)
  }

  err = db.Ping()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error: failed to ping the DB.\n")
    os.Exit(1)
  }

  fmt.Println("DB connection OK.")
}