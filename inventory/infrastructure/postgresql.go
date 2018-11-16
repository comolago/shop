// infrastructure implementation details
package infrastructure

import (
   "os"
   "strconv"
   "database/sql"
   "fmt"
   "github.com/comolago/shop/inventory/domain"
   _ "github.com/lib/pq"
)

// Constants
const (
   dbhost = "DBHOST"
   dbport = "DBPORT"
   dbuser = "DBUSER"
   dbpass = "DBPASS"
   dbname = "DBNAME"
)

// Attributes
type PostgresqlDb struct {
   conn *sql.DB
   host       string
   port       int
   user       string
   pass       string
   dbname     string
}

// Load configuration from environment variables
func (pg *PostgresqlDb) config() *domain.ErrHandler {
   var port string
   ok := false
   pg.host, ok = os.LookupEnv(dbhost)
   if !ok {
      return &domain.ErrHandler{2, "func (pg PostgresqlDb)", "config", ""}
   }
   port, ok = os.LookupEnv(dbport)
   if !ok {
      return &domain.ErrHandler{3, "func (pg PostgresqlDb)", "config", ""}
   }
   pg.port, _ = strconv.Atoi(port)
   pg.user, ok = os.LookupEnv(dbuser)
   if !ok {
      return &domain.ErrHandler{4, "func (pg PostgresqlDb)", "config", ""}
   }
   pg.pass, ok = os.LookupEnv(dbpass)
   if !ok {
      return &domain.ErrHandler{5, "func (pg PostgresqlDb)", "config", ""}
   }
   pg.dbname, ok = os.LookupEnv(dbname)
   if !ok {
      return &domain.ErrHandler{6, "func (pg PostgresqlDb)", "config", ""}
   }
   return nil
}

func (pg *PostgresqlDb)initDb() *domain.ErrHandler {
   if pg.conn == nil {
      return &domain.ErrHandler{7, "func (pg PostgresqlDb)", "initDb()", ""}
   }
   if _, err := pg.conn.Exec("CREATE TABLE IF NOT EXISTS inventory(id integer NOT NULL, name varchar(200) NOT NULL,quantity integer NOT NULL DEFAULT 0,CONSTRAINT inventory_pk PRIMARY KEY (id));"); err != nil {
      return &domain.ErrHandler{1, "func (pg PostgresqlDb)", "initDb(()", err.Error()}
   }
   /*if _, err := pg.conn.Exec("CREATE TABLE IF NOT EXISTS users(id serial NOT NULL, username TEXT NOT NULL, password TEXT NOT NULL);"); err != nil {
      return &domain.ErrHandler{1, "func (pg PostgresqlDb)", "initDb(()", err.Error()}
   }*/
   return nil
}

func (pg *PostgresqlDb)AddItem(item domain.Item) *domain.ErrHandler {
   if pg.conn == nil {
      return &domain.ErrHandler{7, "func (pg PostgresqlDb)", "AddItem(item domain.Item)", ""}
   }
   if _, err := pg.conn.Exec("INSERT INTO inventory (id,name,quantity) VALUES ($1,$2, $3);", item.Id, item.Name, item.Quantity); err != nil {
      return &domain.ErrHandler{1, "func (pg PostgresqlDb)", "initDb(()", err.Error()}
   }
   return nil
}

// Open connection to database
func (pg *PostgresqlDb)Open() *domain.ErrHandler {
   connErr:= pg.config() 
   if connErr != nil {
      return connErr
   }
   var err error
   psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", pg.host, pg.port, pg.user, pg.pass, pg.dbname)
   pg.conn, err = sql.Open("postgres", psqlInfo)
   if err != nil {
      return &domain.ErrHandler{1, "func (pg PostgresqlDb)", "Open", err.Error()}
   }
   err = pg.conn.Ping()
   if err != nil {
      return &domain.ErrHandler{1, "func (pg PostgresqlDb)", "Open", err.Error()}
   }
   connErr = pg.initDb()
   if connErr != nil {
      return connErr
   }
   return nil
}

// Retrieve an Item by its id from database
func (pg *PostgresqlDb)GetItemById(id int, item *domain.Item) *domain.ErrHandler {
   if pg.conn == nil {
      return &domain.ErrHandler{7, "func (pg PostgresqlDb)", "getItem(inventory *domain.Inventory)", ""}
   }
   recordset := pg.conn.QueryRow("SELECT id, name, quantity FROM inventory WHERE id=$1;", id)
   var err domain.ErrHandler
   switch err := recordset.Scan(&item.Id, &item.Name, &item.Quantity); err {
      case sql.ErrNoRows:
         return &domain.ErrHandler{13, "func (pg PostgresqlDb)", "getItem(inventory *domain.Inventory)", ""}
      case nil:
         return nil
   }
   return &domain.ErrHandler{1, "func (pg PostgresqlDb)", "getItem(inventory *domain.Inventory)", err.Error()}
}

/*
// Retrieve an Item by its id from database
func (pg *PostgresqlDb)GetInventory(inventory *domain.Inventory) *domain.ErrHandler {
   if pg.conn == nil {
      return &domain.ErrHandler{7, "func (pg PostgresqlDb)", "GetInventory(inventory *domain.Inventory)", ""}
   }
   recordset, err := pg.conn.Query("SELECT id, name, quantity FROM inventory;")
   if err != nil {
      return &domain.ErrHandler{1, "func (pg PostgresqlDb)", "GetInventory(inventory *domain.Inventory)", err.Error()}
   }
   defer recordset.Close()
   var cnt :=0
   for recordset.Next() {
      inventory.items.new()
      err = recordset.Scan(
         &inventory.item[cnt].Id,
         &inventory.item[cnt].Name,
         &inventory.item[cnt].Quantity,
      )
      if err != nil {
         return &domain.ErrHandler{1, "func (pg PostgresqlDb)", "GetInventory(inventory *domain.Inventory)", err.Error()}
      }
      cnt++
   }
   err = recordset.Err()
   if err != nil {
      return &domain.ErrHandler{1, "func (pg PostgresqlDb)", "GetInventory(inventory *domain.Inventory)", err.Error()}
   }
   return nil
}*/

func (pg *PostgresqlDb)DelItemById(id int) *domain.ErrHandler {
   if pg.conn == nil {
      return &domain.ErrHandler{7, "func (pg PostgresqlDb)", "DelItemById(id int)", ""}
   }
   if _, err := pg.conn.Exec("DELETE FROM inventory WHERE id=$1;", id); err != nil {
      return &domain.ErrHandler{1, "func (pg PostgresqlDb)", "DelItemById(id int)", err.Error()}
   }
   return nil
}

// Retrieve an Item by its id from database
func (pg *PostgresqlDb)AuthenticateUser(username string, password string) (int, *domain.ErrHandler) {
   if pg.conn == nil {
      return -1, &domain.ErrHandler{7, "func (pg PostgresqlDb)", "AuthenticateUser(username string, password string)", ""}
   }
   recordset := pg.conn.QueryRow("select id from users where username=$1 and password=crypt($2,password);", username, password)

   var id int
   var err domain.ErrHandler
   switch err := recordset.Scan(&id); err {
      case sql.ErrNoRows:
         return -1, &domain.ErrHandler{13, "func (pg PostgresqlDb)", "AuthenticateUser(username string, password string)", ""}
      case nil:
         return id, nil
   }
   return -1, &domain.ErrHandler{1, "func (pg PostgresqlDb)", "AuthenticateUser(username string, password string)", err.Error()}
}

