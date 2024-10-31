package repeater

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/MirekKrassilnikov/todo_list_api_sql/config"
)

func StringToTime(dateString string, layout string) (time.Time, error) {
	parsedDate, err := time.Parse(layout, dateString)
	if err != nil {
		return time.Time{}, err
	}
	return parsedDate, nil
}

func NextDate(now string, date string, repeat string) (string, error) {
	nowTimeTime, err := StringToTime(now, config.Layout)
	if err != nil {
		return "", fmt.Errorf("invalid date: %s", now)
	}
	startDateTimeTime, err := StringToTime(date, config.Layout)
	if err != nil {
		return "", fmt.Errorf("invalid date: %s", date)
	}
	codeAndNumber := strings.Split(repeat, " ")
	if len(codeAndNumber) == 0 || (codeAndNumber[0] != "y" && codeAndNumber[0] != "d") {
		return "", fmt.Errorf("invalid repeat code: %s", repeat)
	}
	var nextTimeString string

	if codeAndNumber[0] == "y" {

		for {
			nextTime := startDateTimeTime.AddDate(1, 0, 0)
			if nextTime.After(nowTimeTime) {
				nextTimeString = nextTime.Format("20060102")
				break
			}
			// Обновляем startDateTimeTime для следующей итерации
			startDateTimeTime = nextTime
		}
		return nextTimeString, nil

	}

	if codeAndNumber[0] == "d" {
		if len(codeAndNumber) != 2 {
			return "", fmt.Errorf("invalid day repeat format: %s", repeat)
		}

		i, err := strconv.Atoi(codeAndNumber[1])
		if err != nil {
			return "", fmt.Errorf("error converting string to int: %s", repeat)
		}
		if i > 400 {
			return "", nil
		}
		/*if codeAndNumber[1] == "1" {
			nextTime := time.Now().Format("20060102")
			return nextTime, nil
		}*/

		for {
			nextTime := startDateTimeTime.AddDate(0, 0, i)
			if nextTime.After(nowTimeTime) || nextTime.Equal(nowTimeTime) {
				nextTimeString = nextTime.Format("20060102")
				break
			}
			// Обновляем startDateTimeTime для следующей итерации
			startDateTimeTime = nextTime
		}
		return nextTimeString, nil
	}

	return "", fmt.Errorf("unknown repeat code: %s", codeAndNumber[0])
}
