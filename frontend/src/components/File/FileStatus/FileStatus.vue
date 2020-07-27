<template>
  <div class="file-status">
    <div class="file-status__speed">
      <span class="file-status__speed-item">
        Down: {{ downloadBandwidth | prettyBytes(1) }}/s</span
      >
      <span class="file-status__speed-item"
        >Up: {{ uploadBandwidth | prettyBytes(1) }}/s</span
      >
    </div>

    <canvas
      class="file-status__progress"
      ref="canvas"
      width="400"
      height="6"
    ></canvas>
  </div>
</template>

<style lang="scss">
@import "./FileStatus.scss";
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
  data() {
    return {
      downloadBandwidth: 0,
      uploadBandwidth: 0,
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
        this.drawProgress(newEvent.ChunkMap);
      }
    },
  },
  mounted() {
    this.getChunkMap();
  },
  methods: {
    getChunkMap() {
      window.backend.getFileChunkMap(this.file.FileHash, 400).then((bits) => {
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
        ctx.lineTo(i, 6);
        ctx.closePath();
        ctx.stroke();
      });
    },
  },
};
</script>
