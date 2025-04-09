import * as ecarts from "echarts"
import type { Ref } from 'vue';
function renderChart(chart: Ref<null, null>, x: Ref<string[], string[]>, y: Ref<number[], number[]>) {
  const options = {
    // X 轴配置
    xAxis: {
      type: 'category',
      show: true,
      boundaryGap: true,
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        show: true,
        color: '#666',
        fontSize: 12,
        interval: 0
      },
      data: x.value
    },
    yAxis: {
      type: 'value',
      show: false,
    },
    grid: {
      left: '0%',
      right: '0%',
      top: '20%',
      bottom: '25%',
      containLabel: false
    },
    series: [
      {
        name: 'Temperature',
        type: 'line',
        smooth: true,
        symbol: 'emptyCircle',
        symbolSize: 10,
        data: y.value,
        lineStyle: {
          color: '#5470C6',
          width: 2
        },
        itemStyle: {
          color: '#5470C6',
        },
        label: {
          show: true,
          position: 'top',
          formatter: '{c}°',
          color: '#333',
          fontSize: 14,
          distance: 8
        }
      }
    ]
  };
  const myChart = ecarts.init(chart.value);
  myChart.setOption(options);
}
function renderChart2(chart: Ref<null, null>, x: Ref<string[], string[]>, y: Ref<number[], number[]>) {
  const options = {
    // X 轴配置
    xAxis: {
      type: 'category',
      show: true,
      boundaryGap: true,
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        show: true,
        color: '#666',
        fontSize: 12,
        interval: 0
      },
      data: x.value
    },
    yAxis: {
      type: 'value',
      show: false,
    },
    grid: {
      left: '0%',
      right: '0%',
      top: '20%',
      bottom: '15%',
      containLabel: false
    },
    series: [
      {
        name: 'Temperature',
        type: 'line',
        smooth: true,
        symbol: 'emptyCircle',
        symbolSize: 10,
        data: y.value,
        lineStyle: {
          color: '#5470C6',
          width: 2
        },
        itemStyle: {
          color: '#5470C6',
        },
        label: {
          show: true,
          position: 'top',
          formatter: '{c}%',
          color: '#333',
          fontSize: 14,
          distance: 8
        }
      }
    ]
  };
  const myChart = ecarts.init(chart.value);
  myChart.setOption(options);
}
export {
  renderChart,
  renderChart2,
}