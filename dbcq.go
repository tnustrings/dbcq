package dbcq
import (
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
  _ "github.com/microsoft/go-mssqldb"
  "strings"
  "strconv"
    "os"
    "errors"
    "gopkg.in/yaml.v3"
  "fmt"
  "net/url"
)
type DB struct {
  db *sql.DB
  info Info
}
type Info struct {
  Target string
  Type string      `yaml:"type"`
  Database string  `yaml:"database"`
  User string  `yaml:"user"`
  Password string  `yaml:"password"`
  Host string    `yaml:"host"`
  Port int      `yaml:"port"`
}
type Conf map[string]Info
func Open(target string) (*DB, error) { 
  db := new(DB)
  var err error
  db.info, err = info(target)
  if err != nil { return nil, err }
  db.db, err = sql.Open(db.info.Type, connstr(db.info))
  if err != nil { return nil, err }
  return db, nil
}
func (d *DB) Close() {
  d.db.Close()
}
//Query
func (d *DB) Qfad(query string/*, args ...string*/) ([]map[string]interface{}, error) {
  rows, err := d.db.Query(query) //, args)
  //defer rows.Close()
  if err != nil { return nil, err }
  colnames, err := rows.Columns()
  if err != nil { return nil, err }
  cols := make([]interface{}, len(colnames))
  colptrs := make([]interface{}, len(colnames))
  for i := 0; i < len(colnames); i++ {
    colptrs[i] = &cols[i]
  }
  out := make([]map[string]interface{}, 0)
  for rows.Next() {
    err = rows.Scan(colptrs...)
    if err != nil { return nil, err }
    mymap := make(map[string]interface{})
    
    for i, val := range cols {
      mymap[strings.ToLower(colnames[i])] = val
    }
    for k, val := range mymap {
      switch val.(type) {
      case []uint8:
        s := string(val.([]byte))
        num, err := strconv.Atoi(s)
        if err != nil { return nil, err }
        mymap[k] = num
      }
    }
    out = append(out, mymap)
  }
  rows.Close()
  return out, nil
}
func (d *DB) Info() Info {
  return d.info
}

func info(target string) (Info, error) {
  b, err := os.ReadFile("/home/max/.dbc")
  if err != nil { return Info{}, err }
  cnf := Conf{}
  err = yaml.Unmarshal(b, &cnf)
  if err != nil { return Info{}, err }
  for targ, info := range cnf {
    //fmt.Println("t: ", targ)
    if targ == target {
      info.Target = targ
      return info, nil
    }
  }
  return Info{}, errors.New(fmt.Sprintf("target %s not found.", target))
}
func connstr(info Info) string {
  if info.Type == "sqlite3" {
    return info.Database
  } else if info.Type == "mssql" {
    vals := url.Values{}
    vals.Add("database", info.Database)
    u := &url.URL {
      Scheme: "sqlserver",
      User: url.UserPassword(info.User, info.Password),
      Host: fmt.Sprintf("%s:%d", info.Host, info.Port),
      RawQuery: vals.Encode(),
    }
    return u.String()
  }
    return ""
}
