import Highcharts from 'highcharts'
import HighchartsReact from 'highcharts-react-official'
import HighchartsHeatmap from 'highcharts/modules/heatmap'
import moment from 'moment'
import { DataType } from '../types/data'

interface Props {
  data: DataType
}

if (typeof Highcharts === 'object') {
  HighchartsHeatmap(Highcharts)
}

const HeatMap: React.FC<Props> = ({ data }) => {
  const heatMapData = data.cols
    .filter((a) => {
      const [y, x] = a.X.split('-')

      return !!x && !!y
    })
    .map((a) => {
      const [y, x] = a.X.split('-')

      return [parseInt(x), parseInt(y), a.Y.members]
    })

  const options: Highcharts.Options = {
    chart: {
      type: 'heatmap',
    },
    title: {
      text: '',
    },
    credits: {
      enabled: false,
    },
    legend: {
      enabled: false,
    },
    colorAxis: {
      minColor: '#FFFFFF',
      maxColor: '#2f7ed8',
    },
    tooltip: {
      formatter: function () {
        const day = moment(this.point.y * 60 * 60 * 24 * 1000).format('dddd')
        return (
          day +
          ' @ ' +
          this.point.x +
          ':00 - ' +
          this.point.value.toFixed(0) +
          ' people'
        )
      },
    },
    yAxis: {
      min: 0,
      title: {
        text: 'Day',
      },
      labels: {
        formatter: function () {
          return moment(this.value * 60 * 60 * 24 * 1000).format('dddd')
        },
      },
    },
    xAxis: {
      labels: {
        formatter: function () {
          return this.value.toLocaleString()
        },
      },
    },
    series: [
      {
        name: 'Members',
        //@ts-ignore
        data: heatMapData,
      },
    ],
  }

  return <HighchartsReact highcharts={Highcharts} options={options} />
}

export default HeatMap
