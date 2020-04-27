package dynsql

import (
	"database/sql"
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/gorm"
	"reflect"
	"strings"
	"time"
)

type Model struct {
	Table  string   `json:"table"`
	Fields []string `json:"fields"`
	Wheres []Where  `json:"wheres"`
	Order  string   `json:"order"`
	Limit  uint     `json:"limit"`
	Offset uint     `json:"offset"`
}

type Where struct {
	Query string      `json:"query"`
	Args  interface{} `json:"args"`
	Not   bool        `json:"not"`
	Or    bool        `json:"or"`
}

func Update(rdb *gorm.DB, model *Model) int {
	return 0
}

func Query(rdb *gorm.DB, model *Model) []byte {
	db := rdb.Table(model.Table).Select(strings.Join(model.Fields, ", "))

	for _, w := range model.Wheres {
		var op func(query interface{}, args ...interface{}) *gorm.DB
		if w.Not {
			op = db.Not
		} else if w.Or {
			op = db.Or
		} else {
			op = db.Where
		}
		if w.Args != nil && w.Args != "" {
			db = op(w.Query, w.Args)
		} else {
			db = op(w.Query)
		}
	}

	if model.Order != "" {
		db.Order(model.Order)
	}
	if model.Limit != 0 {
		db.Limit(model.Limit)
	}
	if model.Offset != 0 {
		db.Offset(model.Offset)
	}

	rows, err := db.Rows()
	PanicIfErr(err)

	result, err := RowsToJson(rows)
	PanicIfErr(err)

	bytes, err := json.Marshal(result)
	PanicIfErr(err)

	return bytes
}

func RowsToJson(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.ColumnTypes()
	PanicIfErr(err)
	colNames := make([]string, len(columns))
	for i, column := range columns {
		colNames[i] = strcase.ToLowerCamel(strings.ToLower(column.Name()))
	}
	objects := make([]map[string]interface{}, 0)
	for rows.Next() {
		values := make([]interface{}, len(columns))
		object := map[string]interface{}{}
		for i, column := range columns {
			v := reflect.New(column.ScanType()).Interface()
			switch v.(type) {
			case *sql.RawBytes:
				v = new(JsonNullString)
			case *sql.NullInt64:
				v = new(JsonNullInt64)
			case *sql.NullFloat64:
				v = new(JsonNullFloat64)
			case *mysql.NullTime:
				v = new(*time.Time)
			default:
				v = new(sql.RawBytes)
			}
			object[colNames[i]] = v
			values[i] = v
		}
		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}
		objects = append(objects, object)
	}
	return objects, nil
}

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
