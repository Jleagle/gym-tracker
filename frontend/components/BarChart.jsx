import React from 'react'
import Highcharts from 'highcharts'
import HighchartsReact from 'highcharts-react-official'

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
            data: [], //Gets overridden
        },
        {
            name: 'Capacity %',
            data: [], //Gets overridden
            visible: false,
        }
    ]
}

function BarChart({data}) {

    options.xAxis.categories = Object.keys(data);
    options.series[0].data = Object.values(data).map(a => a.mean_people[0]);
    options.series[1].data = Object.values(data).map(a => a.mean_pcnt[0]);

    return (<HighchartsReact highcharts={Highcharts} options={options}/>);
}

export default BarChart
