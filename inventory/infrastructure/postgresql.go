package infrastructure

import (
   "os"
   "strconv"
   "database/sql"
   "fmt"
   "github.com/comolago/shop/inventory/domain"
   _ "github.com/lib/pq"
)

const (
   dbhost = "DBHOST"
   dbport = "DBPORT"
   dbuser = "DBUSER"
   dbpass = "DBPASS"
   dbname = "DBNAME"
)

type PostgresqlDb struct {
   conn *sql.DB
   host       string
   port       int
   user       string
   pass       string
   dbname     string
}

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

/*CREATE TABLE inventory(id integer NOT NULL,name varchar(200) NOT NULL,quantity integer NOT NULL DEFAULT 0,CONSTRAINT inventory_pk PRIMARY KEY (id));
INSERT INTO inventory VALUES (1,'Fedora Red', 5);

*/
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
      panic(err)
   }
   return nil
}

func (pg *PostgresqlDb)GetItemById(id int, item *domain.Item) *domain.ErrHandler {
   if pg.conn == nil {
      return &domain.ErrHandler{7, "func (pg PostgresqlDb)", "getItem(inventory *domain.Inventory)", ""}
   }
   recordset, err := pg.conn.Query("SELECT id, name FROM inventory WHERE id=$1;", id)
   if err != nil {
      return &domain.ErrHandler{1, "func (pg PostgresqlDb)", "getItem(inventory *domain.Inventory)", err.Error()}
   }
   defer recordset.Close()
   for recordset.Next() {
      err = recordset.Scan(
         &item.Id,
         &item.Name,
      )
      if err != nil {
         return &domain.ErrHandler{1, "func (pg PostgresqlDb)", "getItem(inventory *domain.Inventory)", err.Error()}
      }
   }
   err = recordset.Err()
   if err != nil {
      return &domain.ErrHandler{1, "func (pg PostgresqlDb)", "getItem(inventory *domain.Inventory)", err.Error()}
   }
   return nil
}
