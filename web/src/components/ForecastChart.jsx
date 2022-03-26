import { createEffect, onMount, mergeProps } from 'solid-js'
import styles from './ForecastChart.module.css'
import * as Echarts from 'echarts'
import data from '../data/forecast1'

function ForecastChart(props) {
    const [farms, prepared] = prepareData(data)

    console.log('Farms :', farms)
    console.log('Data  :', prepared)

    const p = mergeProps({ title: 'Forecast - Yield Farming' }, props)

    let ref
    let chart

    onMount(() => {
        chart = Echarts.init(ref, 'dark')
    })

    createEffect(() => {
        const datasets = []
        const seriesList = []

        farms.forEach((name) => {
            const names = [`${name} Farm`, `${name} HODL`, `${name} Only A`]

            names.forEach((name) => {
                const datasetId = 'dataset_' + name

                datasets.push({
                    id: datasetId,
                    fromDatasetId: 'dataset_raw',
                    transform: {
                        type: 'filter',
                        config: {
                            and: [{ dimension: 'Name', '=': name }],
                        },
                    },
                })

                seriesList.push({
                    type: 'line',
                    datasetId: datasetId,
                    showSymbol: false,
                    name,
                    endLabel: {
                        show: true,
                        formatter: function (params) {
                            let n = params.value[1]
                            if (n.length < 20) {
                                n += ' '.repeat(20 - n.length)
                            }
                            let v = params.value[2]
                            return `${n}: $${roundTo2(v).toLocaleString()}`
                        },
                    },
                    labelLayout: {
                        moveOverlap: 'shiftY',
                    },
                    emphasis: {
                        focus: 'series',
                    },
                    encode: {
                        x: 'Date',
                        y: 'Value',
                        label: ['Name', 'Value'],
                        itemName: 'Date',
                        tooltip: ['Value'],
                    },
                    tooltip: {
                        valueFormatter: (value) =>
                            `$${roundTo2(value).toLocaleString()}`,
                        textStyle: {
                            fontFamily: 'Source Code Pro',
                        },
                    },
                })
            })
        })

        const options = {
            animationDuration: 3000,
            dataset: [
                {
                    id: 'dataset_raw',
                    source: prepared,
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
                name: 'Total Value (USD)',
            },
            grid: {
                right: 250,
            },
            series: seriesList,
            textStyle: {
                fontFamily: 'Source Code Pro',
            },
        }

        chart.setOption(options)
    })

    return (
        <div class={styles.ChartContainer}>
            <div class={styles.Chart} ref={ref} />
        </div>
    )
}

function prepareData(data) {
    const names = {}
    const fieldsMap = {}

    // Assume first row of data contains field names: ['Date', 'Name', 'Value', 'HODL', ...]
    data[0].forEach((f, i) => (fieldsMap[f] = i))

    const nameIdx = fieldsMap['Name']
    const dateIdx = fieldsMap['Date']
    const valueIdx = fieldsMap['Value']
    const holdIdx = fieldsMap['HODL']
    const onlyAIdx = fieldsMap['Only A']
    const aprIdx = fieldsMap['APR']

    data.forEach((v, i) => {
        if (i === 0) {
            return
        }

        const name = v[nameIdx]
        if (!names[name]) {
            names[name] = []
        }

        const date = v[dateIdx]
        const value = v[valueIdx]
        const hold = v[holdIdx]
        const onlyA = v[onlyAIdx]
        const apr = v[aprIdx]

        names[name].push(
            [date, name + ' Farm', value, apr], // Farm value row.
            [date, name + ' HODL', hold, 0.0], // HOLD value row.
            [date, name + ' Only A', onlyA, 0.0] // Only A value row.
        )
    })

    const uniqueNames = Object.keys(names)
    const out = uniqueNames.reduce(
        (acc, key) => {
            acc.push(...names[key])
            return acc
        },
        [['Date', 'Name', 'Value', 'APR']]
    )
    return [uniqueNames, out]
}

function findUniqueFields(fieldName, data) {
    const fieldsMap = {}
    const uniqueFieldsNames = data.reduce((acc, v, i) => {
        if (i === 0) {
            v.forEach((f, i) => (fieldsMap[f] = i))
            if (fieldsMap[fieldName] == null) {
                console.error(
                    `could not find '${fieldName} among field names:`,
                    fieldsMap
                )
            }
        } else {
            acc.add(v[fieldsMap[fieldName]])
        }
        return acc
    }, new Set())
    return [[...uniqueFieldsNames], fieldsMap]
}

function roundTo2(num) {
    return Math.round((num + Number.EPSILON) * 100) / 100
}

export default ForecastChart
