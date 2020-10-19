<template>
  <div class="file-chunks">
    <div class="file-chunks__percent">{{ progress.toFixed(2) }}%</div>
    <canvas
      class="file-chunks__progress"
      ref="canvas"
      width="120"
      height="20"
    ></canvas>
  </div>
</template>

<style lang="scss">
@import "./FileChunks.scss";
</style>

<script>
import { mapState } from "vuex";

export default {
  props: {
    file: {
      type: Object,
      default: () => {},
    },
  },
  data: () => {
    return {
      downloadBandwidth: 0,
      uploadBandwidth: 0,
      progress: 0,
    };
  },
  computed: {
    ...mapState("downloadEvents", ["downloadEvent"]),
  },
  watch: {
    downloadEvent(newEvent) {
      const { FileHash } = this.file;
      if (FileHash === newEvent.FileHash) {
        this.downloadBandwidth = newEvent.DownloadBandwidth;
        this.uploadBandwidth = newEvent.UploadBandwidth;
        this.progress = newEvent.Progress * 100;
        this.drawProgress(newEvent.ChunkMap);
      }
    },
    progress(x) {
      if (x === 100) {
        this.$store.dispatch("files/fetchLocalFiles");
      }
    },
  },
  mounted() {
    this.getChunkMap();
    this.progress = !this.file.IsPaused && !this.file.IsDownloading ? 100 : 0;
  },
  methods: {
    getChunkMap() {
      window.backend.getFileChunkMap(this.file.FileHash, 118).then((bits) => {
        this.drawProgress(bits);
      });
    },
    drawProgress(bits) {
      const canvas = this.$refs.canvas;
      const ctx = canvas.getContext("2d");
      const colours = ["#ebebeb", "#02d2b3"];

      const bitmap = `${bits}`.split("");

      bitmap.forEach((val, i) => {
        ctx.beginPath();
        ctx.strokeStyle = colours[parseFloat(val)];
        ctx.lineWidth = 1;
        ctx.moveTo(i, 0);
        ctx.lineTo(i, 20);
        ctx.closePath();
        ctx.stroke();
      });
    },
  },
};
</script>
