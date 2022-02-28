package timeseries

import (
	"math"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	dateFormat = "2006-01-02"
)

type Series []ValueAt

func (ts Series) At(t time.Time) (price ValueAt, found bool) {
	return ts.AtDate(t.Format(dateFormat))
}

func (ts Series) AtDate(date string) (price ValueAt, found bool) {
	for _, p := range ts {
		if p.Date() == date {
			price = p
			found = true
			break
		}
	}
	return
}

type ValueAt struct {
	TS int64   `json:"ts"`
	V  float64 `json:"v"`
}

func (v ValueAt) Time() time.Time {
	return time.UnixMilli(int64(v.TS))
}

func (v ValueAt) Date() string {
	return v.Time().Format(dateFormat)
}

func (ts Series) Print() {
	pr := message.NewPrinter(language.English)

	for _, v := range ts {
		pr.Printf("[%s]  --  %.04f\n", v.Date(), v.V)
	}
}

func FromTuples(tuples [][]interface{}) Series {
	var ts Series
	for _, tuple := range tuples {
		t := ValueAt{
			TS: int64(tuple[0].(float64)),
			V:  tuple[1].(float64),
		}
		ts = append(ts, t)
	}
	return ts
}

func DiffDays(dateA, dateB string) (days int) {
	a, err1 := time.Parse(dateFormat, dateA)
	b, err2 := time.Parse(dateFormat, dateB)
	if err1 == nil && err2 == nil {
		days = int(math.Abs(b.Sub(a).Hours() / 24))
	}
	return
}

func ToTime(date string) time.Time {
	t, _ := time.Parse(dateFormat, date)
	return t
}

func ToDate(t time.Time) (date string) {
	return t.Format(dateFormat)
}
