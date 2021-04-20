import Highcharts from 'highcharts'
import HighchartsReact from 'highcharts-react-official'
import HighchartsHeatmap from 'highcharts/modules/heatmap'
import { DataType } from '../types/data'

interface Props {
  data: DataType
}

if (typeof Highcharts === 'object') {
  HighchartsHeatmap(Highcharts)
}

const days = ['', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];

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
        return this.point.value.toFixed(1) + ' people on '
            + days[this.point.y] + 's at '
            + this.point.x + ':00 ' +' on average'
      },
    },
    yAxis: {
      min: 1,
      max: 7,
      reversed: true,
      gridLineWidth: 0,
      title: {
        text: '',
      },
      labels: {
        formatter: function () {
          return days[this.value]
        },
      },
    },
    xAxis: {
      gridLineWidth: 0,
      type: 'category',
      labels: {
        step:1,
        formatter: function () {
          return this.value.toLocaleString() + ':00'
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
