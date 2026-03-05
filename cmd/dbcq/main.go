package main
import (
  "log"
  "flag"
  "fmt"
  "github.com/tnustrings/dbcq"
  "encoding/json"
)
func main() {
  flag.Parse()
  args := flag.Args()
  target := args[0]
  query := args[1]
  db, err := dbcq.Open(target)
  defer db.Close()
  if err != nil { log.Fatal("error: ", err) }
  out, err := db.Qfad(query)
  if err != nil { fmt.Println(err) }
  b, err := json.Marshal(out)
  if err != nil { log.Fatal(err) }
  fmt.Println(string(b))
}
