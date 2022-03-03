package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/anrid/traderbot/pkg/timeseries"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/spf13/pflag"
)

func main() {
	path := pflag.StringP("path", "p", "", "path to output dir (required, e.g. /mnt/c/Users/whatever/)")
	pflag.Parse()

	if *path == "" {
		pflag.PrintDefaults()
		os.Exit(-1)
	}

	page := components.NewPage()
	page.AddCharts(
		lineBase(),
	)

	file := filepath.Join(*path, "line.html")
	fmt.Printf("writing chart %s\n", file)

	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}

func generateLineItems(max int) []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < max; i++ {
		items = append(items, opts.LineData{Value: rand.Intn(300)})
	}
	return items
}

func lineBase() *charts.Line {
	var dates []string
	today := time.Now().UTC()
	start := today.Add(-90 * 24 * time.Hour)
	for ; start.Before(today); start = start.Add(24 * time.Hour) {
		d := timeseries.ToDate(start)
		dates = append(dates, d)
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{Title: "Title", Subtitle: "Sub."}),
	)

	line.SetXAxis(dates).
		AddSeries("Category A", generateLineItems(len(dates))).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	return line
}
