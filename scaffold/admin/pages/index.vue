<template>
  <div class="anal-chart">
    <div class="level u-plr">
      <div class="level-left">
        <div class="level-item">
          <button class="button " @click="prevMonth">
            <span class="icon">
              <i class="fa fa-chevron-left"></i>
            </span>
            &nbsp;Prev Month
          </button>
        </div>
      </div>
      <div class="level-item">
        <h3 class="title is-5">
          Unique Viewers for {{monthName}}
          {{year}}
        </h3>
      </div>
      <div class="level-right">
        <div class="level-item">
          <button class="button " @click="nextMonth">
            Next Month&nbsp;
            <span class="icon">
              <i class="fa fa-chevron-right"></i>
            </span>
          </button>
        </div>
      </div>
    </div>
    <div class="u-rel">
      <b-loading :is-full-page="false" :active="isLoading"></b-loading>
      <chartist v-if="!isLoading"
        ratio="ct-major-second"
        type="Line"
        :data="chartData"
        :options="chartOptions" >
      </chartist>
    </div>
  </div>
</template>

<script>
// ... imports

import Chartist from 'vue-bulma-chartist'

export default {
  // COMPONENT
  // ______________________________________
  name: 'Home',
  layout: 'admin',
  components: {
    Chartist
  },
  props: {},
  computed: {
    monthName () {
      if (this.month === 1) return 'January'
      if (this.month === 2) return 'February'
      if (this.month === 3) return 'March'
      if (this.month === 4) return 'April'
      if (this.month === 5) return 'May'
      if (this.month === 6) return 'June'
      if (this.month === 7) return 'July'
      if (this.month === 8) return 'August'
      if (this.month === 9) return 'September'
      if (this.month === 10) return 'October'
      if (this.month === 11) return 'November'
      if (this.month === 12) return 'December'
    }
  },
  methods: {
    prevMonth () {
      var month = this.month - 1
      if (month === 0) {
        month = 12
        this.year--
      }
      this.month = month
      this.load()
    },
    nextMonth () {
      var month = this.month + 1
      if (month > 12) {
        month = 1
        this.year++
      }
      this.month = month
      this.load()
    },
    load () {
      this.isLoading = true
      this.$axios.get(`api/v1/views/${this.month}/${this.year}`).then(({data}) => {
        this.isLoading = false
        this.chartData = data
      })
    }
  },
  watch: {},
  data () {
    return {
      month: new Date().getMonth() === 0 ? 12 : new Date().getMonth(), // last month
      year: new Date().getMonth() === 0 ? new Date().getFullYear() - 1 : new Date().getFullYear(),
      isLoading: true,
      chartData: {
        labels: [],
        series: []
      },
      chartOptions: {
        lineSmooth: false,
        height: 500,
        axisY: {
          labelInterpolationFnc: function (value, index, data) {
            if ((value + '').indexOf('.') >= 0) {
              return ''
            }
            return value
          }
        }
      }
    }
  },

  // LIFECYCLE METHODS
  // ______________________________________
  beforeCreate () {
  },
  created () {
    // this.$router.replace({ name: 'pages', params: { ID: 0 } })
    this.load()
  },
  beforeMount () {
  },
  mounted () {
  },
  beforeUpdate () {
  },
  updated () {
  },
  beforeDestroy () {
  },
  destroyed () {
  }
}
</script>

<style lang="scss">
@import "~public/scss/_variables";

.anal-chart {
  .ct-series-a .ct-point, .ct-series-a .ct-line, .ct-series-a .ct-bar, .ct-series-a .ct-slice-donut {
    stroke: $blue;
  }
  .ct-series-b .ct-point, .ct-series-b .ct-line, .ct-series-b .ct-bar, .ct-series-b .ct-slice-donut {
    stroke: $red;
  }
  .ct-series-c .ct-point, .ct-series-c .ct-line, .ct-series-c .ct-bar, .ct-series-c .ct-slice-donut {
    stroke: $yellow;
  }
}
</style>
