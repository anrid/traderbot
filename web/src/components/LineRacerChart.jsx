import { createEffect, onMount, mergeProps } from 'solid-js'
import styles from './LineRacerChart.module.css'
import * as Echarts from 'echarts'
import data from '../data/pop'

function Hello(props) {
    const p = mergeProps(
        { title: 'Income in various countries since 1950' },
        props
    )

    let ref
    let chart

    onMount(() => {
        chart = Echarts.init(ref, 'dark')
    })

    createEffect(() => {
        const countries = [
            'Finland',
            'France',
            'Germany',
            'Iceland',
            'Norway',
            'Poland',
            'Russia',
            'United Kingdom',
        ]

        const datasets = []
        const seriesList = []

        Echarts.util.each(countries, function (country) {
            var datasetId = 'dataset_' + country

            datasets.push({
                id: datasetId,
                fromDatasetId: 'dataset_raw',
                transform: {
                    type: 'filter',
                    config: {
                        and: [
                            { dimension: 'Year', gte: 1950 },
                            { dimension: 'Country', '=': country },
                        ],
                    },
                },
            })

            seriesList.push({
                type: 'line',
                datasetId: datasetId,
                showSymbol: false,
                name: country,
                endLabel: {
                    show: true,
                    formatter: function (params) {
                        return (
                            params.value[3] +
                            ': ' +
                            params.value[0].toLocaleString()
                        )
                    },
                },
                labelLayout: {
                    moveOverlap: 'shiftY',
                },
                emphasis: {
                    focus: 'series',
                },
                encode: {
                    x: 'Year',
                    y: 'Income',
                    label: ['Country', 'Income'],
                    itemName: 'Year',
                    tooltip: ['Income'],
                },
            })
        })

        const options = {
            animationDuration: 3000,
            dataset: [
                {
                    id: 'dataset_raw',
                    source: data,
                },
                ...datasets,
            ],
            title: {
                text: p.title,
            },
            tooltip: {
                order: 'valueDesc',
                trigger: 'axis',
            },
            xAxis: {
                type: 'category',
                nameLocation: 'middle',
            },
            yAxis: {
                name: 'Income',
            },
            grid: {
                right: 140,
            },
            series: seriesList,
        }

        chart.setOption(options)
    })

    return (
        <div class={styles.ChartContainer}>
            <div class={styles.Chart} ref={ref} />
        </div>
    )
}

export default Hello
