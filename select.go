package sqlutils

import (
	"database/sql"
	"fmt"
	"reflect"
)

func DbSelect(db *sql.DB, query string, args ...any) ([]map[string]string, []string) {
	rows, err := db.Query(query, args...)
	defer rowResurect(rows)
	if err != nil {
		panic(err)
	}
	return rowsToMaps(rows)
}

func rowsToMaps(rows *sql.Rows) ([]map[string]string, []string) {
	var datas []map[string]string
	cols, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		scans := make([]interface{}, len(cols))
		row := make(map[string]string)
		for i := range scans {
			scans[i] = &scans[i]
		}
		if err = rows.Scan(scans...); err != nil {
			panic(err)
		}
		for i, v := range scans {
			value := ""
			if v != nil {
				if reflect.TypeOf(v).String() == "string" || reflect.TypeOf(v).String() == "[]uint8" {
					value = fmt.Sprintf("%s", v)
				} else {
					value = fmt.Sprintf("%v", v)
				}
			}
			row[cols[i]] = value
		}
		datas = append(datas, row)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}
	return datas, cols
}

func rowResurect(rows *sql.Rows) {
	NormalError()
	if err := rows.Close(); err != nil {
		panic(err)
	}
}
