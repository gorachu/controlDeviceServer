package sqlite

import (
	"database/sql"
	"fmt"
	"strconv"
)

func (s *Storage) InsertShield(row map[string]interface{}) error {
	const query = `
INSERT INTO shields (
    number_in_list,
    shipping_date,
    has_controller,
    has_input,
    has_output,
    has_analyzer,
    has_lcd,
    customer,
    inspector,
    install_address,
    has_sim,
    phone_num
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)

`
	getString := func(key string) string {
		if v, ok := row[key]; ok && v != nil {
			return fmt.Sprintf("%v", v)
		}
		return ""
	}

	parseInt := func(key string) sql.NullInt64 {
		v := getString(key)
		if v == "" {
			return sql.NullInt64{Valid: false}
		}
		if num, err := strconv.ParseInt(v, 10, 64); err == nil {
			return sql.NullInt64{Int64: num, Valid: true}
		}
		return sql.NullInt64{Valid: false}
	}

	hasSim := getString("Наличие SIM") != "Без SIM"

	_, err := s.db.Exec(query,
		parseInt("№"),
		getString("Дата отгрузки"),
		getString("№ PLC"),
		getString("№ Input"),
		getString("№ Output"),
		getString("№ Анализатор тока"),
		getString("№ Экран"),
		getString("Заказчик"),
		getString("ОТК"),
		getString("Адрес установки"),
		hasSim,
		getString("Номер телефона"),
	)
	return err
}
