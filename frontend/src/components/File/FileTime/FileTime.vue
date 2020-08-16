<template>
  <div class="file-time">
    <div class="file-time__title">
      <template
        v-if="!file.IsDownloading && !file.IsUploading && !file.IsPaused"
      >
        Ready
      </template>
      <template v-else-if="progress === 100 || file.IsUploading">
        Seeding
      </template>
      <template v-else-if="file.IsPaused">
        Paused
      </template>
      <template v-else>
        {{ [seconds, "seconds"] | duration("humanize", true) }}
      </template>
    </div>
    <div class="file-time__percent">{{ progress.toFixed(2) }}%</div>
  </div>
</template>

<style lang="scss">
@import "./FileTime.scss";
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
      downloadBandwidth: 0,
      seconds: 0,
    };
  },
  mounted() {
    this.progress = !this.file.IsPaused && !this.file.IsDownloading ? 100 : 0;
  },
  computed: {
    ...mapState("downloadEvents", ["downloadEvent"]),
  },
  watch: {
    downloadEvent(newEvent) {
      if (this.file.FileHash === newEvent.FileHash) {
        this.seconds =
          (this.file.FileSize - this.file.FileSize * newEvent.Progress) /
          newEvent.DownloadBandwidth;
        this.downloadBandwidth = newEvent.DownloadBandwidth;
        this.progress = newEvent.Progress * 100;
      }
    },
  },

  methods: {},
};
</script>
