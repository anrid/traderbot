package coingecko

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEMAIndicator(t *testing.T) {
	r := require.New(t)

	now := time.Now().UTC()
	_24h := 24 * time.Hour

	day1 := now.Add(_24h * -4).UnixMilli()
	day2 := now.Add(_24h * -3).UnixMilli()
	day3 := now.Add(_24h * -2).UnixMilli()
	day4 := now.Add(_24h * -1).UnixMilli()
	day5 := now.UnixMilli()

	prices := NewTimeSeries([][]interface{}{
		{float64(day1), 1.0},
		{float64(day2), 2.0},
		{float64(day3), 3.0},
		{float64(day4), 4.0},
		{float64(day5), 5.0},
	})

	i := NewEMAIndicator(3, prices)
	// Dump(i)

	r.Equal(4.5, i.V[day5])
	r.Equal(4.25, i.V[day4])
	r.Equal(3.625, i.V[day3])
	r.Equal(0.0, i.V[day2])
	r.Equal(0.0, i.V[day1])
}
