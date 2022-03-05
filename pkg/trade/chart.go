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

	fontFamily := "Source Code Pro"

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme:  types.ThemeVintage,
			Width:  "1000px",
			Height: "700px",
		}),
		charts.WithTitleOpts(opts.Title{
			Title:    title,
			Subtitle: subtitle,
			TitleStyle: &opts.TextStyle{
				FontFamily: fontFamily,
			},
			SubtitleStyle: &opts.TextStyle{
				FontFamily: fontFamily,
			},
		}),
		charts.WithLegendOpts(opts.Legend{
			Show:   true,
			Bottom: "1px",
			TextStyle: &opts.TextStyle{
				FontSize:   12,
				FontFamily: fontFamily,
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Type: "value",
			AxisLabel: &opts.AxisLabel{
				Show:         true,
				ShowMaxLabel: true,
				ShowMinLabel: true,
			},
			SplitLine: &opts.SplitLine{
				Show: true,
				LineStyle: &opts.LineStyle{
					Type: "dotted",
				},
			},
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Date",
			AxisLabel: &opts.AxisLabel{
				Show:         true,
				ShowMaxLabel: true,
				ShowMinLabel: true,
			},
			SplitLine: &opts.SplitLine{
				Show: true,
				LineStyle: &opts.LineStyle{
					Type: "dotted",
				},
			},
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show: true,
		}),
	)

	latest := len(farm.TotalValueHistory) - 1
	farmV := farm.TotalValueHistory[latest]
	farmPL := ((farmV / farm.InitialInvestment) - 1) * 100
	hodlV := farm.TotalValueHODLHistory[latest]
	hodlPL := ((hodlV / farm.InitialInvestment) - 1) * 100
	onlyAV := farm.TotalValueHODLOnlyAHistory[latest]
	onlyAPL := ((onlyAV / farm.InitialInvestment) - 1) * 100
	onlyBV := farm.TotalValueHODLOnlyBHistory[latest]
	onlyBPL := ((onlyBV / farm.InitialInvestment) - 1) * 100

	series1 := "Farm"
	series2 := "HODL"
	series3 := pr.Sprintf("Only %s", strings.ToUpper(farm.A.Symbol))
	series4 := pr.Sprintf("Only %s", strings.ToUpper(farm.B.Symbol))

	line.SetXAxis(farm.ChangeHistory).
		AddSeries(pr.Sprintf("%s: %.f (%.f%%)", series1, farmV, farmPL), toLineData(farm.TotalValueHistory)).
		AddSeries(pr.Sprintf("%s: %.f (%.f%%)", series2, hodlV, hodlPL), toLineData(farm.TotalValueHODLHistory)).
		AddSeries(pr.Sprintf("%s: %.f (%.f%%)", series3, onlyAV, onlyAPL), toLineData(farm.TotalValueHODLOnlyAHistory)).
		AddSeries(pr.Sprintf("%s: %.f (%.f%%)", series4, onlyBV, onlyBPL), toLineData(farm.TotalValueHODLOnlyBHistory))

	// line.SetSeriesOptions(
	// 	charts.WithLineChartOpts(
	// 		opts.LineChart{Smooth: true}),
	// 	),
	// )

	markPointNames := []string{series1, series2, series3, series4}

	for i := 0; i < len(line.MultiSeries); i++ {
		name := markPointNames[i]

		line.MultiSeries[i].MarkPoints = &opts.MarkPoints{
			Data: []interface{}{
				opts.MarkLineNameTypeItem{Name: name + " max", Type: "max"},
				opts.MarkLineNameTypeItem{Name: name + " min", Type: "min"},
			},
			MarkPointStyle: opts.MarkPointStyle{
				Label: &opts.Label{
					Show:      true,
					Position:  "top",
					Formatter: "{b}: {c}",
				},
				SymbolSize: 10.0,
				Symbol:     []string{"diamond"},
			},
		}
	}

	page := components.NewPage()
	page.AddCharts(line).SetLayout(components.PageFlexLayout)

	file := filepath.Join(path, filename)
	fmt.Printf("writing chart %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		return errors.Wrapf(err, "could not write chart to file %s", file)
	}
	page.Render(io.MultiWriter(f))

	return nil
}
