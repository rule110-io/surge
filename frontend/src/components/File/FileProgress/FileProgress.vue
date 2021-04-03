<template>
  <div class="progress-bar">
    <div
      class="progress-bar__value"
      :style="{ width: progressFormatted }"
    ></div>
    <div class="progress-bar__text">
      {{ progressFormatted }}
    </div>
  </div>
</template>

<style lang="scss">
@import "./FileProgress.scss";
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
  computed: {
    ...mapState("globalBandwidth", ["statusBundle"]),
    progressFormatted() {
      return `${this.progress.toFixed(2)}%`;
    },
  },
  watch: {
    statusBundle(newEvent) {
      const { FileHash } = this.file;
      const newFileHash = this._.find(newEvent, { FileHash });
      const isNewFileHash = !this._.isEmpty(newFileHash);
      if (isNewFileHash) {
        this.shared = newFileHash.ChunksShared / newFileHash.NumChunks;
        this.progress = newFileHash.Progress * 100;
        this.drawProgress(newFileHash.ChunkMap);
      }
    },
    progress(x) {
      if (x === 100) {
        this.$store.dispatch("files/fetchLocalFiles");
      }
    },
  },
  data: () => {
    return {
      progress: 0,
    };
  },
  mounted() {
    this.progress = this.file.Progress * 100;
  },
  methods: {},
};
</script>
