<template>
  <LineChart
    class="bandwidth-chart"
    :chart-data="chartData"
    :options="options"
  ></LineChart>
</template>

<style lang="scss">
@import "./BandwidthChart.scss";
</style>

<script>
import { mapState } from "vuex";
import LineChart from "@/components/Charts/Bandwidth.js";
import "@taeuk-gang/chartjs-plugin-streaming";

export default {
  props: {
    file: {
      type: Object,
      default: () => {},
    },
  },
  components: {
    LineChart,
  },
  computed: {
    ...mapState("globalBandwidth", ["statusBundle"]),
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
            backgroundColor: "#2CF2FF",
            borderColor: "#2CF2FF",
            borderWidth: 2,
            fill: false,
          },
          {
            label: "up",
            data: [0],
            backgroundColor: "#FB49C0",
            borderColor: "#FB49C0",
            borderWidth: 2,
            fill: false,
          },
        ],
      },
      options: {
        maintainAspectRatio: false,
        responsive: true,
        // animation: {
        //   duration: 0,
        // },
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
              display: true,
              stacked: true,
              grid: {
                display: true,
              },
              gridLines: {
                borderDash: [4.8, 7.2],
                borderDashOffset: 2,
                color: "rgba(255,255,255,0.2)",
                drawBorder: false,
              },
              ticks: {
                beginAtZero: true,
                min: 0,
                stepSize: 0.25,
                fontColor: "#7B7D82",
                callback: function (label) {
                  return `${label} Mb/s `;
                },
              },
            },
          ],
          xAxes: [
            {
              display: false,
              type: "realtime",
              gridLines: {
                zeroLineColor: "transparent",
              },
            },
          ],
        },
        plugins: {
          streaming: {
            onRefresh: (chart) => {
              if (!this.file) return;

              const { FileHash } = this.file;
              const newFileHash = this._.find(this.statusBundle, { FileHash });
              const isNewFileHash = !this._.isEmpty(newFileHash);

              if (isNewFileHash) {
                chart.data.labels.push(Date.now());
                chart.data.datasets[0].data.push(
                  (newFileHash.DownloadBandwidth / 1000000).toFixed(2)
                );
                chart.data.datasets[1].data.push(
                  (newFileHash.UploadBandwidth / 1000000).toFixed(2)
                );
              }
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
