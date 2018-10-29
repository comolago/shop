package domain

import (
   "os"
   "strconv"
   "database/sql"
   "fmt"
   _ "github.com/lib/pq"
)




const (
   dbhost = "DBHOST"
   dbport = "DBPORT"
   dbuser = "DBUSER"
   dbpass = "DBPASS"
   dbname = "DBNAME"
)

// Store struct definition
type PostgresqlDb struct {
   conn *sql.DB
   host       string
   port       int
   user       string
   pass       string
   dbname     string
}

func (pg *PostgresqlDb) config() *ErrHandler {
   ok := false
   pg.host, ok = os.LookupEnv(dbhost)
   if !ok {
      return &ErrHandler{2, "func (pg PostgresqlDb)", "config", ""}
   }
   port := ""
   port, ok = os.LookupEnv(dbport)
   if !ok {
      return &ErrHandler{3, "func (pg PostgresqlDb)", "config", ""}
   }
   pg.port, _ = strconv.Atoi(port)
   fmt.Println(fmt.Sprintf("port:%d",pg.port))
   pg.user, ok = os.LookupEnv(dbuser)
   if !ok {
      return &ErrHandler{4, "func (pg PostgresqlDb)", "config", ""}
   }
   pg.pass, ok = os.LookupEnv(dbpass)
   if !ok {
      return &ErrHandler{5, "func (pg PostgresqlDb)", "config", ""}
   }
   pg.dbname, ok = os.LookupEnv(dbname)
   if !ok {
      return &ErrHandler{6, "func (pg PostgresqlDb)", "config", ""}
   }
   return nil
}

func (pg *PostgresqlDb)Open() *ErrHandler {

   fmt.Println("Open")
   connErr:= pg.config() 
   fmt.Println(fmt.Sprintf("Port:%d",pg.port))
   if connErr != nil {
      return connErr
   }
   var err error
   psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", pg.host, pg.port, pg.user, pg.pass, pg.dbname)
   pg.conn, err = sql.Open("postgres", psqlInfo)
   if err != nil {
      return &ErrHandler{1, "func (pg PostgresqlDb)", "Open", err.Error()}
   }
   err = pg.conn.Ping()
   if err != nil {
      return &ErrHandler{1, "func (pg PostgresqlDb)", "Open", err.Error()}
      panic(err)
   }
   //fmt.Println("Successfully connected!")
   return nil
}

func (pg *PostgresqlDb)getItem(inventory *Inventory) *ErrHandler {
   /*if pg.conn == nil {
      return &ErrHandler{7, "func (pg PostgresqlDb)", "getItem(inventory *Inventory)", ""}
   }
   recordset, err := pg.conn.Query(`
      SELECT
         id, name,
      FROM items
      ORDER BY id ASC`)
   if err != nil {
      return &ErrHandler{1, "func (pg PostgresqlDb)", "getItem(inventory *Inventory)", err.Error()}
   }
   defer recordset.Close()
   for recordset.Next() {
      item := Item{}
      err = recordset.Scan(
         &item.Id,
         &item.Name,
      )
      if err != nil {
         return &ErrHandler{1, "func (pg PostgresqlDb)", "getItem(inventory *Inventory)", err.Error()}
      }
      inventory.Items = append(inventory.Items, item)
   }
   err = recordset.Err()
   if err != nil {
      return &ErrHandler{1, "func (pg PostgresqlDb)", "getItem(inventory *Inventory)", err.Error()}
   }*/
   return nil
}
