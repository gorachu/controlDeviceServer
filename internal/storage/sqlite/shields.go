package sqlite

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

func (s *Storage) InsertShield(row map[string]interface{}) error {
	const query = `
INSERT INTO shields (
    number_in_list,
    shipping_date,
    customer,
    inspector,
    install_address,
    has_sim,
    phone_num
) VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(number_in_list) DO UPDATE SET
    shipping_date = excluded.shipping_date,
    customer = excluded.customer,
    inspector = excluded.inspector,
    install_address = excluded.install_address,
    has_sim = excluded.has_sim,
	phone_num = excluded.phone_num;
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
		getString("Заказчик"),
		getString("ОТК"),
		getString("Адрес установки"),
		hasSim,
		getString("Номер телефона"),
	)
	if err != nil {
		return fmt.Errorf("ошибка вставки в shields: %w", err)
	}
	moduleFields := []struct {
		key   string
		table string
	}{
		{"№ PLC", "controllers"},
		{"№ Input", "inputmodules"},
		{"№ Анализатор тока", "current_analyzers"},
		{"№ Экран", "lcds"},
	}

	for _, field := range moduleFields {
		nums, err := ParseCommaSeparatedInts(getString(field.key))
		if err != nil {
			return fmt.Errorf("ошибка парсинга поля %s: %w", field.key, err)
		}

		for _, num := range nums {
			_, err := s.db.Exec(
				fmt.Sprintf("UPDATE %s SET in_shield = ? WHERE number_in_list = ?", field.table),
				parseInt("№ PLC"), num,
			)
			if err != nil {
				return fmt.Errorf("ошибка обновления in_shield для %s (№%d): %w", field.table, num, err)
			}
		}
	}

	return nil
}
func ParseCommaSeparatedInts(input string) ([]int, error) {
	parts := strings.Split(input, ",")

	var result []int
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		num, err := strconv.Atoi(trimmed)
		if err != nil {
			return nil, fmt.Errorf("не удалось преобразовать '%s' в число: %w", trimmed, err)
		}
		result = append(result, num)
	}
	return result, nil
}
