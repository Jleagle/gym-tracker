import React from 'react'
import Highcharts from 'highcharts'
import HighchartsReact from 'highcharts-react-official'

function BarChart({data}) {

    const options = {
        chart: {
            type: 'column',
        },
        title: {
            text: '',
        },
        credits: {
            enabled: false,
        },
        xAxis: {
            crosshair: true,
            categories: Object.keys(data),
        },
        yAxis: [
            {
                min: 0,
                title: {
                    text: 'Members'
                },
                labels: {
                    formatter: function () {
                        return this.value.toLocaleString();
                    },
                },
            },
            {
                min: 0,
                max: 100,
                title: {
                    text: 'Capacity'
                },
                labels: {
                    formatter: function () {
                        return this.value + ' %';
                    },
                },
                opposite: true,
            }
        ],
        series: [
            {
                name: 'Members',
                data: Object.values(data).map(a => a.members),
            },
            {
                name: 'Capacity',
                data: Object.values(data).map(a => a.percent),
                visible: false,
            }
        ]
    }

    return (<HighchartsReact highcharts={Highcharts} options={options}/>);
}

export default BarChart
