<template>
  <div
    class="file-card"
    :style="{ background: `${activeColor}` }"
    @click="openFile(file)"
  >
    <div class="file-card__header">
      <feather class="file-card__icon" type="folder"></feather>
      <div class="file-card__size">{{ file.FileSize | prettyBytes(1) }}</div>
    </div>
    <div class="file-card__footer">
      <div class="file-card__title text_wrap_none">{{ file.FileName }}</div>
      <div class="file-card__progress">
        Finished — {{ progress.toFixed(2) }}%
      </div>
    </div>
  </div>
</template>

<style lang="scss">
@import "./FileCard.scss";
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
  computed: {
    ...mapState("downloadEvents", ["downloadEvent"]),
  },
  watch: {
    downloadEvent(newEvent) {
      if (this.file.FileHash === newEvent.FileHash) {
        this.progress = newEvent.Progress * 100;
      }
    },
  },
  data: () => {
    return {
      colors: ["#FEC606", "#2CC990", "#03bf7b", "#8870FF"],
      progress: 0,
    };
  },
  created() {
    this.progress = !this.file.IsPaused && !this.file.IsDownloading ? 100 : 0;
    this.activeColor = this.getRandomColor();
  },
  methods: {
    openFile(file) {
      const { FileHash } = file;
      window.backend.openFile(FileHash).then(() => {});
    },
    getRandomColor() {
      const colors = this.colors;
      return colors[Math.floor(Math.random() * colors.length)];
    },
  },
};
</script>
