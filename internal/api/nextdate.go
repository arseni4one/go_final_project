package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату выполнения задачи
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("правило повторения не указано")
	}

	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %s", err)
	}

	parts := strings.Split(repeat, " ")
	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат правила d")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("неверное количество дней (должно быть 1-400)")
		}
		for {
			date = date.AddDate(0, 0, days)
			if date.After(now) {
				return date.Format("20060102"), nil
			}
		}

	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				return date.Format("20060102"), nil
			}
		}

	default:
		return "", fmt.Errorf("неподдерживаемый формат правила: %s", parts[0])
	}
}
