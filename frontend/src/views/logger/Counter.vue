<template>
  <div :id="id" :class="className" :style="{height:height,width:width}" />
</template>

<script>
import echarts from 'echarts/lib/echarts'
import 'echarts/lib/chart/line'
import 'echarts/lib/chart/bar'
import 'echarts/lib/component/title'
import 'echarts/lib/component/tooltip'
import 'echarts/lib/component/legend'
import resize from '@/utils/resize'

export default {
  mixins: [resize],
  props: {
    className: {
      type: String,
      default: 'chart'
    },
    id: {
      type: String,
      default: 'chart'
    },
    width: {
      type: String,
      default: '90%'
    },
    height: {
      type: String,
      default: '400px'
    },
    chartdata: {
      type: Object,
      default: () => {}
    }
  },
  data() {
    return {
      chart: null,
      colors: {
        default: ['#c23531', '#2f4554', '#61a0a8', '#d48265', '#91c7ae', '#749f83', '#ca8622', '#bda29a', '#6e7074', '#546570', '#c4ccd3'],
        light: ['#37A2DA', '#32C5E9', '#67E0E3', '#9FE6B8', '#FFDB5C', '#ff9f7f', '#fb7293', '#E062AE', '#E690D1', '#e7bcf3', '#9d96f5', '#8378EA', '#96BFFF'],
        dark: ['#dd6b66', '#759aa0', '#e69d87', '#8dc1a9', '#ea7e53', '#eedd78', '#73a373', '#73b9bc', '#7289ab', '#91ca8c', '#f49f42']
      }
    }
  },
  computed: {
    xAxis() {
      return this.chartdata.xAxis
    },
    legend() {
      return this.chartdata.legend
    },
    series() {
      return this.chartdata.series
    }
  },
  watch: {
    chartdata() {
      this.setOption()
    }
  },
  mounted() {
    this.initChart()
  },
  beforeDestroy() {
    if (!this.chart) return
    this.chart.dispose()
    this.chart = null
  },
  methods: {
    initChart() {
      this.chart = echarts.init(document.getElementById(this.id))
      this.chart.on('click', (params) => {
        this.$emit('click-event', {
          seriesName: params.seriesName,
          name: params.name,
          value: params.value
        })
      })
    },
    setOption() {
      this.chart.clear()
      this.chart.setOption({
        title: { orient: 'vertical', x: 'center', text: '' },
        tooltip: { trigger: 'axis' },
        grid: { containLabel: true, x: 96, y: 20, y2: 20 },
        color: '#9FE6B8',
        xAxis: { type: 'category', boundaryGap: true, data: this.xAxis },
        yAxis: { type: 'value', min: 0, logBase: 10 },
        series: this.series
      })
    },
    resizeChart() {
      if (this.chart) {
        this.chart.resize()
      }
    }
  }
}
</script>
