import React from 'react'
import Highcharts from 'highcharts'
import HighchartsReact from 'highcharts-react-official'
import moment from 'moment';

function LineChart({data}) {

    const options = {
        chart: {
            type: 'spline',
        },
        title: {
            text: ''
        },
        credits: {
            enabled: false,
        },
        legend: {
            verticalAlign: 'bottom',
        },
        xAxis: {
            crosshair: true,
            categories: data.cols.map(a => moment(a.X * 1000).format("HH:mm")),
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
                data: data.cols.map(a => a.Y.members),
            },
            {
                name: 'Capacity',
                data: data.cols.map(a => a.Y.percent),
                visible: false,
            }
        ],
    };

    return (<HighchartsReact highcharts={Highcharts} options={options}/>);
}

export default LineChart
