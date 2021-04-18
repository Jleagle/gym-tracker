import React from 'react'
import Highcharts from 'highcharts'
import HighchartsReact from 'highcharts-react-official'
import ordinalSuffix from "../helpers/ordinal";

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
            categories: function () {

                return data.cols.map(function (a) {

                    if (a.X === '') {
                        return '';
                    }

                    switch (data.group) {
                        case 'yearDay':
                        case 'monthDay':
                            return ordinalSuffix(a.X);
                        case 'weekDay':
                            const days = ['', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];
                            return days[a.X];
                        case 'weekHour':
                            return '';
                        case 'hour':
                            return a.X;
                        default:
                            return '';
                    }
                });
            }(),
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
                gridLineWidth: 0,
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
        ]
    }

    return (<HighchartsReact highcharts={Highcharts} options={options}/>);
}

export default BarChart
