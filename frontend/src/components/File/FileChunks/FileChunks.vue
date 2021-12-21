<template>
  <div class="file-chunks">
    <canvas class="file-chunks__progress" ref="canvas"></canvas>
  </div>
</template>

<style lang="scss">
@import "./FileChunks.scss";
</style>

<script>
import { mapState } from "vuex";

export default {
  components: {},
  props: {
    file: {
      type: Object,
      default: () => {},
    },
  },
  data: () => {
    return {
      progress: 0,
      canvasWidth: 0,
      canvasHeight: 31,
    };
  },
  computed: {
    ...mapState("globalBandwidth", ["statusBundle"]),
  },
  watch: {
    statusBundle(newEvent) {
      const { FileHash } = this.file;
      const newFileHash = this._.find(newEvent, { FileHash });
      const isNewFileHash = !this._.isEmpty(newFileHash);

      if (isNewFileHash) {
        this.getChunkMap();
      }
    },
    "file.FileHash": {
      handler(newHash) {
        console.log(newHash);
        if (!newHash) return;

        this.getChunkMap();
      },
    },
  },
  mounted() {
    this.getChunkMap();
  },
  methods: {
    getChunkMap() {
      const canvasWidth = this.$refs.canvas.clientWidth;
      this.canvasWidth = canvasWidth;

      window.go.surge.MiddlewareFunctions.GetFileChunkMap(
        this.file.FileHash,
        canvasWidth
      ).then((bits) => {
        this.drawProgress(bits);
      });
    },
    drawProgress(bits) {
      const canvas = this.$refs.canvas;
      canvas.width = this.canvasWidth;
      canvas.height = this.canvasHeight;
      const ctx = canvas.getContext("2d");
      const colours = ["#242425", "#086099", "#0BA193"];

      const bitmap = `${bits}`.split("");

      bitmap.forEach((val, i) => {
        ctx.beginPath();
        ctx.strokeStyle = colours[parseFloat(val)];
        ctx.lineWidth = 1;
        ctx.moveTo(i, 0);
        ctx.lineTo(i, this.canvasHeight);
        ctx.closePath();
        ctx.stroke();
      });
    },
  },
};
</script>
