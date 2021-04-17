import React from 'react'
import Highcharts from 'highcharts'
import HighchartsReact from 'highcharts-react-official'
import moment from "moment";

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
            categories: data.map(a => moment(a.X * 1000).format("DD MMM @ HH:mm")),
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
                data: data.map(a => a.Y.members),
            },
            {
                name: 'Capacity',
                data: data.map(a => a.Y.percent),
                visible: false,
            }
        ]
    }

    return (<HighchartsReact highcharts={Highcharts} options={options}/>);
}

export default BarChart
