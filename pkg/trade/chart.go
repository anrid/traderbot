package trade

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func RenderYieldFarmingPerformanceChart(path string, farm *LPFarm) error {
	toLineData := func(vs []float64) (items []opts.LineData) {
		for _, v := range vs {
			items = append(items, opts.LineData{Value: int64(v)})
		}
		return
	}

	pr := message.NewPrinter(language.English)

	title := pr.Sprintf("Yield Farming %s/%s LP  --  [%s - %s]",
		strings.ToUpper(farm.A.Symbol),
		strings.ToUpper(farm.B.Symbol),
		farm.ChangeHistory[0],
		farm.ChangeHistory[len(farm.ChangeHistory)-1],
	)
	subtitle := pr.Sprintf("Starting APR: %.f%% , Final APR: %.f%% , Initial Investment: %.f %s",
		farm.APRHistory[0],
		farm.APRHistory[len(farm.APRHistory)-1],
		farm.InitialInvestment,
		strings.ToUpper(string(farm.Currency)),
	)

	filename := strings.ToLower(pr.Sprintf("yield-farming-%s-%s-%s-%s.html",
		farm.A.Symbol,
		farm.B.Symbol,
		farm.ChangeHistory[0],
		farm.ChangeHistory[len(farm.ChangeHistory)-1],
	))

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeMacarons}),
		charts.WithTitleOpts(opts.Title{Title: title, Subtitle: subtitle}),
	)

	line.SetXAxis(farm.ChangeHistory).
		AddSeries("Farm", toLineData(farm.TotalValueHistory)).
		AddSeries("HODL", toLineData(farm.TotalValueHODLHistory)).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{Smooth: true}),
			charts.WithMarkPointNameTypeItemOpts(
				opts.MarkPointNameTypeItem{Name: "Maximum", Type: "max"},
				// opts.MarkPointNameTypeItem{Name: "Average", Type: "average"},
				// opts.MarkPointNameTypeItem{Name: "Minimum", Type: "min"},
			),
			charts.WithMarkPointStyleOpts(
				opts.MarkPointStyle{Label: &opts.Label{Show: true}}),
		)

	page := components.NewPage()
	page.AddCharts(line)

	file := filepath.Join(path, filename)
	fmt.Printf("writing chart %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		return errors.Wrapf(err, "could not write chart to file %s", file)
	}
	page.Render(io.MultiWriter(f))

	return nil
}