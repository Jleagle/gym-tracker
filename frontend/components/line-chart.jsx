import Highcharts from 'highcharts'
import HighchartsReact from 'highcharts-react-official'
import moment from 'moment'

function LineChart({data}) {

  const options = {
    chart: {
      type: 'spline',
      marginTop: 30,
      height: 400,
    },
    title: {
      text: '',
    },
    credits: {
      enabled: false,
    },
    legend: {
      verticalAlign: 'bottom',
    },
    plotOptions: {
      series: {
        lineWidth: 5,
        marker: {
          enabled: false,
        },
      },
    },
    tooltip: {
      formatter: function () {
        switch (this.series.name) {
          case 'Members':
            return '<b>' + this.y + '</b> people at <b>' + this.x + '</b>'
          case 'Capacity':
            return '<b>' + this.y.toFixed(1) + '</b>% full at <b>' + this.x + '</b>'
        }
      },
    },
    xAxis: {
      crosshair: true,
      categories: data.cols.map(a => moment(a.X * 1000).format('HH:mm')),
      labels: {
        step: 1,
        formatter: function () {
          if (this.value.endsWith(':00')) {
            return this.value
          }
          return ''
        },
      },
    },
    yAxis: [
      {
        min: 0,
        title: {
          text: 'Members',
        },
        labels: {
          formatter: function () {
            return this.value.toLocaleString()
          },
        },
        gridLineWidth: 0,
      },
      {
        min: 0,
        max: 100,
        title: {
          text: 'Capacity',
        },
        labels: {
          formatter: function () {
            return this.value + ' %'
          },
        },
        opposite: true,
        visible: false,
      },
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
      },
    ],
  }

  return (
    <div className="chart">
      <HighchartsReact highcharts={Highcharts} options={options}/>
    </div>
  )
}

export default LineChart
