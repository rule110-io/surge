<template>
  <div class="file-time">
    <div class="file-time__title">ETA</div>
    <div class="file-time__status">
      <template v-if="file.IsPaused">
        Paused
      </template>
      <template v-else-if="progress === 100">
        Finished
      </template>
      <template v-else>
        {{ [seconds, "seconds"] | duration("humanize") }}
      </template>
    </div>
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
    progress(x) {
      if (x === 100) {
        this.$store.dispatch("files/fetchLocalFiles");
      }
    },
    downloadEvent(newEvent) {
      if (this.file.FileHash === newEvent.FileHash) {
        this.seconds =
          (this.file.FileSize - this.file.FileSize * newEvent.Progress) /
          newEvent.DownloadBandwidth;
        this.progress = newEvent.Progress * 100;
      }
    },
  },

  methods: {},
};
</script>
