package sqlutils

import (
	"database/sql"
)

func DbSelect(db *sql.DB, query string, args ...any) ([]map[string]interface{}, []string) {
	rows, err := db.Query(query, args...)
	defer rowResurect(rows)
	if err != nil {
		panic(err)
	}
	return rowsToMaps(rows)
}

func rowsToMaps(rows *sql.Rows) ([]map[string]interface{}, []string) {
	var datas []map[string]interface{}
	cols, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		scans := make([]interface{}, len(cols))
		row := make(map[string]interface{})
		for i := range scans {
			scans[i] = &scans[i]
		}
		if err = rows.Scan(scans...); err != nil {
			panic(err)
		}
		for i, v := range scans {
			if v != nil {
				row[cols[i]] = v
			}
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
	if rows != nil {
		if err := rows.Close(); err != nil {
			panic(err)
		}
	}
}
