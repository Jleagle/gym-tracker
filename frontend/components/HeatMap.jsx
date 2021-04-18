import React from 'react'
import Highcharts from 'highcharts'
import HighchartsReact from 'highcharts-react-official'
import HighchartsHeatmap from "highcharts/modules/heatmap";
import moment from "moment";

if (typeof Highcharts === 'object') {
    HighchartsHeatmap(Highcharts);
}

function HeatMap({data}) {

    const options = {
        chart: {
            type: 'heatmap',
        },
        title: {
            text: ''
        },
        credits: {
            enabled: false,
        },
        legend: {
            enabled: false,
        },
        colorAxis: {
            minColor: '#FFFFFF',
            maxColor: '#2f7ed8'
        },
        tooltip: {
            formatter: function () {

                const day = moment(this.point.y * 60 * 60 * 24 * 1000).format("dddd");
                return day + ' @ ' + this.point.x + ':00 - ' + this.point.value.toFixed(0) + ' people';
            }
        },
        yAxis: {
            min: 0,
            title: {
                text: 'Day'
            },
            labels: {
                formatter: function () {
                    return moment(this.value * 60 * 60 * 24 * 1000).format("dddd");
                },
            },
        },
        xAxis: {
            labels: {
                formatter: function () {
                    return this.value.toLocaleString();
                },
            },
        },
        series: {
            name: 'Members',
            data: data.cols.filter(function (a) {
                const [y, x] = a.X.split('-');
                return (Boolean(x) && Boolean(y));
            }).map(function (a) {
                const [y, x] = a.X.split('-');
                return [parseInt(x), parseInt(y), a.Y.members];
            }),
        },
    };

    return (
        <HighchartsReact highcharts={Highcharts} options={options}/>
    );
}

export default HeatMap
