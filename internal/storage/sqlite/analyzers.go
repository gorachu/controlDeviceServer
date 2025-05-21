package sqlite

import (
	"database/sql"
	"fmt"
	"strconv"
)

func (s *Storage) InsertCurrentAnalyzer(row map[string]interface{}) error {
	query := `
		INSERT OR IGNORE INTO current_analyzers (
    		number_in_list,
			type,
			firmware,
			inspector,
			comment
		) VALUES (?, ?, ?, ?, ?);
	`

	get := func(k string) string {
		if v, ok := row[k]; ok && v != nil {
			return fmt.Sprintf("%v", v)
		}
		return ""
	}
	parseInt := func(k string) sql.NullInt64 {
		v := get(k)
		if v == "" {
			return sql.NullInt64{Valid: false}
		}
		if num, err := strconv.ParseInt(v, 10, 64); err == nil {
			return sql.NullInt64{Int64: num, Valid: true}
		}
		return sql.NullInt64{Valid: false}
	}

	_, err := s.db.Exec(query,
		parseInt("№"),
		get("Тип"),
		get("FW"),
		get("ОТК"),
		get("Комментарий"),
	)
	return err
}
