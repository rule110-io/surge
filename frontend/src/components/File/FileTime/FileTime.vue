<template>
  <div class="file-time text_wrap_none">
    <template v-if="file.IsDownloading && !file.IsPaused">
      {{ [seconds, "seconds"] | duration("humanize") }}
    </template>
    <template v-else> ∞ </template>
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
    ...mapState("globalBandwidth", ["statusBundle"]),
  },
  watch: {
    progress(x, oldVal) {
      if (x === 100 && oldVal) {
        this.$store.dispatch("files/fetchLocalFiles");
      }
    },
    statusBundle(newEvent) {
      const { FileHash } = this.file;
      const newFileHash = this._.find(newEvent, { FileHash });
      const isNewFileHash = !this._.isEmpty(newFileHash);

      if (isNewFileHash) {
        this.seconds =
          (this.file.FileSize - this.file.FileSize * newFileHash.Progress) /
          newFileHash.DownloadBandwidth;
        this.progress = newFileHash.Progress * 100;
      }
    },
  },

  methods: {},
};
</script>
