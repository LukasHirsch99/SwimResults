package updatemeet

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

func parseDeadline(s string) (time.Time, error) {
	return time.Parse("02.01.2006 15:04", s)
}


func extractMeetDate(s string) (time.Time, time.Time, error) {
	reFullRange := regexp.MustCompile(`(\d{2})\.-(\d{2})\.(\d{2})\.(\d{4})`)
	reCrossMonth := regexp.MustCompile(`(\d{2})\.(\d{2})\.-(\d{2})\.(\d{2})\.(\d{4})`)
	reSingle := regexp.MustCompile(`(\d{2})\.(\d{2})\.(\d{4})`)

	var startStr string
	var endStr string

	switch {
	// 29.02.-01.03.2020
	case reCrossMonth.MatchString(s):
		match := reCrossMonth.FindStringSubmatch(s)
		startStr = fmt.Sprintf("%s.%s.%s", match[1], match[2], match[5])
		endStr = fmt.Sprintf("%s.%s.%s", match[3], match[4], match[5])

	// 02.-03.10.2020
	case reFullRange.MatchString(s):
		match := reFullRange.FindStringSubmatch(s)
		startStr = fmt.Sprintf("%s.%s.%s", match[1], match[3], match[4])
		endStr = fmt.Sprintf("%s.%s.%s", match[2], match[3], match[4])

	// 03.10.2020
	case reSingle.MatchString(s):
		match := reSingle.FindStringSubmatch(s)
		startStr = fmt.Sprintf("%s.%s.%s", match[1], match[2], match[3])
		endStr = fmt.Sprintf("%s.%s.%s", match[1], match[2], match[3])

	default:
		return time.Time{}, time.Time{}, errors.New("Could not parse dates")
	}

	startDate, err := time.Parse("02.01.2006", startStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("error parsing date: '%s'", startStr)
	}
	endDate, err := time.Parse("02.01.2006", endStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("error parsing date: '%s'", endStr)
	}

	return startDate, endDate, nil
}
