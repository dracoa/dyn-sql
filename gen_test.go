package dynsql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

var db *gorm.DB
var err error

func init() {
	_ = godotenv.Load()
	db, err = gorm.Open("mysql", os.Getenv("CONN_STR"))
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)
	db.DB().Ping()
	log.Println("db init success")
}

func TestParse(t *testing.T) {
	model := &Model{
		Table:  "tkd_item_detail",
		Fields: []string{"ID_ITEM_ID", "ID_DESC", "ID_STOCK_CONTROL_AT", "ID_STOCK_SNAPSHOT"},
		Wheres: []Where{{
			Query: "ID_STOCK_CONTROL_AT IS NOT NULL",
			Args:  "",
			Not:   false,
			Or:    false,
		}},
		Order:  "ID_ITEM_ID",
		Limit:  0,
		Offset: 0,
	}
	result := Query(db, model)
	defer db.Close()
	log.Println(string(result))
}
