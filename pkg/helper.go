package pkg

import "time"

const DATETIME_FORMAT = "2006-01-02 15:04:05"

func DatetimeCanonical(date *time.Time) *string {
	if date == nil {
		return nil
	}

	loc, _ := time.LoadLocation("America/Sao_Paulo")
	d := date.In(loc).Format(DATETIME_FORMAT)
	return &d
}

func TimeFromCanonical(datetime *string) (*time.Time, error) {
	if datetime == nil {
		return nil, nil
	}

	t, err := time.Parse(DATETIME_FORMAT, *datetime)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
