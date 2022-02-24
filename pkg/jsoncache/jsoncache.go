package jsoncache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type InvalidateCachePeriod int

var (
	ErrNotFound = os.ErrNotExist
)

const (
	InvalidateHourly InvalidateCachePeriod = iota + 1
	InvalidateDaily
	InvalidateWeekly
)

func Get(key string, into interface{}, i InvalidateCachePeriod) error {
	k := createKey(key, i)
	return readJSON(k, into)
}

func Set(key string, data interface{}, i InvalidateCachePeriod) error {
	k := createKey(key, i)
	return writeJSON(k, data)
}

func readJSON(key string, into interface{}) error {
	path := filepath.Join(os.TempDir(), key)

	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.Wrapf(err, "could not stat temp file %s", path)
		}
		return ErrNotFound
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "could not read JSON from temp file %s", path)
	}

	err = json.Unmarshal(b, into)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal data")
	}

	return nil
}

func writeJSON(key string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "could not marshal data")
	}

	path := filepath.Join(os.TempDir(), key)

	err = ioutil.WriteFile(path, b, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "could not write JSON to temp file %s", path)
	}

	return nil
}

var (
	wordCharsOnly = regexp.MustCompile(`\W+`)
)

func createKey(input string, i InvalidateCachePeriod) string {
	now := time.Now().UTC()

	var prefix string
	switch i {
	case InvalidateHourly:
		prefix = now.Format("2006-01-02-15")
	case InvalidateDaily:
		prefix = time.Now().Format("2006-01-02")
	case InvalidateWeekly:
		_, week := now.ISOWeek()
		prefix = time.Now().Format("2006-01-") + fmt.Sprintf("week%02d", week)
	}

	return prefix + "-" + wordCharsOnly.ReplaceAllString(strings.ToLower(input), "-")
}
