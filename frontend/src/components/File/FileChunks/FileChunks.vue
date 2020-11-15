<template>
  <div class="file-chunks">
    <div class="file-chunks__title text_wrap_none">
      <template
        v-if="!file.IsDownloading && !file.IsUploading && !file.IsPaused"
      >
        Finished
      </template>
      <template v-else-if="file.IsUploading">
        Seeding
      </template>
      <template v-else-if="file.IsPaused">
        Paused: {{ progress.toFixed(2) }}%
      </template>
      <template v-else> Downloading: {{ progress.toFixed(2) }}% </template>
    </div>
    <canvas
      class="file-chunks__progress"
      ref="canvas"
      width="156"
      height="12"
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
      progress: 0,
    };
  },
  computed: {
    ...mapState("downloadEvents", ["downloadEvent"]),
    baseColor() {
      return document.getElementById("app").classList.contains("dark-theme")
        ? "#ebebeb"
        : "#fcfcfc";
    },
  },
  watch: {
    downloadEvent(newEvent) {
      const { FileHash } = this.file;
      if (FileHash === newEvent.FileHash) {
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
      window.backend.getFileChunkMap(this.file.FileHash, 156).then((bits) => {
        this.drawProgress(bits);
      });
    },
    drawProgress(bits) {
      const canvas = this.$refs.canvas;
      const ctx = canvas.getContext("2d");
      const colours = [this.baseColor, "#5EC1FF", "#02d2b3"];

      const bitmap = `${bits}`.split("");

      bitmap.forEach((val, i) => {
        ctx.beginPath();
        ctx.strokeStyle = colours[parseFloat(val)];
        ctx.lineWidth = 1;
        ctx.moveTo(i, 0);
        ctx.lineTo(i, 12);
        ctx.closePath();
        ctx.stroke();
      });
    },
  },
};
</script>
