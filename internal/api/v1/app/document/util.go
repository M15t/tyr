package document

import (
	"regexp"
	"strconv"
	"time"

	"github.com/araddon/dateparse"
)

func extractNumbers(input string) []int {
	re := regexp.MustCompile("[0-9]+")
	matches := re.FindAllString(input, -1)

	var numbers []int
	for _, match := range matches {
		num, err := strconv.Atoi(match)
		if err == nil {
			numbers = append(numbers, num)
		}
	}
	return numbers
}

func addDaysToDate(dateString string, daysToAdd int) (string, error) {
	// Parse the date string into a time.Time object
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return "", err
	}

	// Add the specified number of days
	newDate := date.AddDate(0, 0, daysToAdd)

	// Format the new date as a string
	newDateString := newDate.Format("2006-01-02")

	return newDateString, nil
}

func parseStringToDate(dateString string) string {
	t, _ := dateparse.ParseAny(dateString)

	return t.Format("2006-01-02")
}
