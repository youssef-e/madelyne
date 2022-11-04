package isotime

import (
	"strings"
	"time"
)

type IsoTime struct {
	t     time.Time
	valid bool
}

const layout = "2006-01-02T15:04:05+0000"

func New(year, month, day int) IsoTime {
	if year == 0 || month == 0 {
		return IsoTime{}
	}
	if year != 0 && month != 0 && day == 0 {
		day = 10
	}
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return IsoTime{t: t, valid: true}

}

func Now() IsoTime {
	return IsoTime{t: time.Now(), valid: true}

}

func (mt *IsoTime) UnmarshalJSON(data []byte) error {
	stringData := strings.Trim(string(data), "\"")
	if stringData == "" {
		return nil
	}
	t, err := time.Parse(layout, stringData)
	if err != nil {
		return err
	}
	mt.t = t
	mt.valid = true
	return nil
}

func (mt IsoTime) MarshalJSON() ([]byte, error) {
	if mt.IsValid() {
		return []byte(`"` + mt.t.Format(layout) + `"`), nil
	}
	return []byte(`""`), nil
}

func (mt *IsoTime) IsValid() bool {
	return mt.valid
}
