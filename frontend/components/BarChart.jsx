import React from 'react'
import Highcharts from 'highcharts'
import HighchartsReact from 'highcharts-react-official'
import ordinalSuffix from "../helpers/ordinal";
import borderRadius from 'highcharts-border-radius';

if (typeof Highcharts === 'object') {
    borderRadius(Highcharts);
}

const days = ['', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];

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
        plotOptions: {
            column: {
                borderRadiusTopLeft: 5,
                borderRadiusTopRight: 5,
            },
        },
        tooltip: {
            formatter: function () {

                let ret = '<b>' + this.y.toFixed(1) + '</b>';

                switch (this.series.name) {
                    case 'Members':
                        ret += ' people';
                        break;
                    case 'Capacity':
                        ret += '% full';
                        break;
                }

                switch (data.group) {
                    case 'yearDay':
                        ret += ' on the ' + this.x + ' day'
                        break;
                    case 'monthDay':
                        ret += ' on the ' + this.x
                        break;
                    case 'weekDay':
                        ret += ' on ' + this.x + 's'
                        break;
                    case 'hour':
                        ret += ' at ' + this.x
                        break;
                }

                return ret + ' on average';
            }
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
                            return days[a.X];
                        case 'weekHour':
                            return '';
                        case 'hour':
                            return a.X + ':00';
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
                visible: false
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
