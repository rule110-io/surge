<template>
  <line-chart
    class="bandwidth-chart"
    :chart-data="chartData"
    :options="options"
  ></line-chart>
</template>

<style lang="scss">
@import "./BandwidthChart.scss";
</style>

<script>
import { mapState } from "vuex";
import LineChart from "@/components/Charts/Bandwidth.js";

export default {
  components: {
    LineChart,
  },
  computed: {
    ...mapState("globalBandwidth", ["totalDown", "totalUp"]),
  },
  watch: {},
  data() {
    return {
      chartData: {
        labels: [1, 2],
        datasets: [
          {
            label: "down",
            data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            backgroundColor: "rgba(44, 201, 144, 0.2)",
            borderColor: "rgba(44, 201, 144, 1)",
            borderWidth: 0,
          },
          {
            label: "up",
            data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            backgroundColor: "rgba(214, 54, 73, 0.1)",
            borderColor: "rgba(214, 54, 73, 1",
            borderWidth: 0,
          },
        ],
      },
      options: {
        animation: {
          duration: 0,
        },
        responsive: false,
        legend: {
          display: false,
        },
        elements: {
          point: {
            radius: 0,
          },
        },
        tooltips: {
          enabled: false,
        },
        scales: {
          yAxes: [
            {
              display: false,
              ticks: {
                beginAtZero: true,
                min: 0,
              },
            },
          ],
          xAxes: [
            {
              display: false,
            },
          ],
        },
      },
    };
  },
  mounted() {
    this.interval = setInterval(this.updateData, 1000);
  },
  destroyed() {
    clearInterval(this.interval);
  },
  methods: {
    updateData() {
      const chartData = Object.assign({}, this.chartData);
      chartData.datasets[0].data.shift();
      chartData.datasets[0].data.push(this.totalDown);

      chartData.datasets[1].data.shift();
      chartData.datasets[1].data.push(this.totalUp);

      this.chartData = chartData;
    },
  },
};
</script>
