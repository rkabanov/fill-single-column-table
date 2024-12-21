package main

import (
  "bytes"
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

  // read file into an array
//  values make([]string, 200)
  data, err := os.ReadFile(*args.file)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error: failed to read file %v.\n", *args.file)
    os.Exit(1)
  }
 
  values := bytes.Split(data, []byte("\n"))
  fmt.Printf("Read %v lines.\n", len(values))

  // insert
  query := fmt.Sprintf("insert into %v(%v) values ($1)", *args.table, *args.column)
  fmt.Println("Query:", query)

  count := 0
  for i := range(values) {
    _, err = db.Exec(query, values[i])
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error: failed to insert value #%v: %v.\n", i, values[i])
      panic(err)
    }
    count++
    if i % 100 == 0 {
      fmt.Printf("%v...\n", i);
    }
  }

  fmt.Printf("Inserted %v records.\n", count)
}