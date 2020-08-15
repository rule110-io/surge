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
import "chartjs-plugin-streaming";

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
            data: [0],
            backgroundColor: "rgba(44, 201, 144, 0.3)",
            borderColor: "rgba(44, 201, 144, 1)",
            borderWidth: 1,
          },
          {
            label: "up",
            data: [0],
            backgroundColor: "rgba(91, 152, 220, 0.2)",
            borderColor: "rgba(91, 152, 220, 1)",
            borderWidth: 1,
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
              type: "realtime",
            },
          ],
        },
        plugins: {
          streaming: {
            onRefresh: (chart) => {
              chart.data.labels.push(Date.now());
              chart.data.datasets[0].data.push(this.totalDown);
              chart.data.datasets[1].data.push(this.totalUp);
            },
            delay: 2000,
          },
        },
      },
    };
  },
  mounted() {},
  destroyed() {},
  methods: {},
};
</script>
