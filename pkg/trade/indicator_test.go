package trade

import (
	"testing"
	"time"

	"github.com/anrid/traderbot/pkg/timeseries"
	"github.com/stretchr/testify/require"
)

func TestEMAIndicator(t *testing.T) {
	r := require.New(t)

	now := time.Now().UTC()
	_24h := 24 * time.Hour

	day1 := now.UnixMilli()
	day2 := now.Add(_24h * 1).UnixMilli()
	day3 := now.Add(_24h * 2).UnixMilli()
	day4 := now.Add(_24h * 3).UnixMilli()
	day5 := now.Add(_24h * 4).UnixMilli()

	prices := timeseries.FromTuples([][]interface{}{
		{float64(day1), 1.0},
		{float64(day2), 2.0},
		{float64(day3), 3.0},
		{float64(day4), 4.0},
		{float64(day5), 5.0},
	})

	i := NewEMAIndicator(3, prices)
	// Dump(i)

	r.Equal(0.0, i.ForTimestamp(day1))
	r.Equal(0.0, i.ForTimestamp(day2))
	r.Equal(0.0, i.ForTimestamp(day3))
	r.Equal(3.0, i.ForTimestamp(day4))
	r.Equal(4.0, i.ForTimestamp(day5))
}
