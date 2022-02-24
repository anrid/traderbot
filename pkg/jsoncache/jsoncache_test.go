package jsoncache

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateKey(t *testing.T) {
	r := require.New(t)

	now := time.Now().UTC()
	thisHour := now.Format("2006-01-02-15")
	thisDate := now.Format("2006-01-02")
	_, week := now.ISOWeek()
	thisWeek := now.Format("2006-01-") + fmt.Sprintf("week%02d", week)

	r.Equal(createKey("ABC", InvalidateHourly), thisHour+"-abc")
	r.Equal(createKey("@ABC@", InvalidateDaily), thisDate+"--abc-")
	r.Equal(createKey("a/B/c", InvalidateWeekly), thisWeek+"-a-b-c")
}

func TestReadWriteJSON(t *testing.T) {
	r := require.New(t)

	type curry struct {
		MassaMun string `json:"massa_mun"`
		Gai      string `json:"gai"`
	}

	data := curry{
		"curry",
		"chicken",
	}

	key1 := "testing/testing/123"
	key2 := "testing/testing/does-not-exist"

	r.NoError(Set(key1, data, InvalidateWeekly))

	data2 := curry{}

	r.NoError(Get(key1, &data2, InvalidateWeekly))
	r.Equal(data.MassaMun, data2.MassaMun)
	r.Equal(data.Gai, data2.Gai)

	r.Equal(ErrNotFound, Get(key2, &data2, InvalidateWeekly))
}
