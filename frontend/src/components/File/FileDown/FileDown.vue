<template>
  <div class="file-down text_wrap_none">
    <template v-if="downloadBandwidth > 0">
      {{ downloadBandwidth | prettyBytes(0) }}/s
    </template>
    <template v-else> - </template>
  </div>
</template>

<style lang="scss">
@import "./FileDown.scss";
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
    };
  },
  mounted() {},
  computed: {
    ...mapState("globalBandwidth", ["statusBundle"]),
  },
  watch: {
    statusBundle(newEvent) {
      if (!this.file) return;

      const { FileHash } = this.file;
      const newFileHash = this._.find(newEvent, { FileHash });
      const isNewFileHash = !this._.isEmpty(newFileHash);

      if (isNewFileHash) {
        this.downloadBandwidth = newFileHash.DownloadBandwidth;
      }
    },
  },

  methods: {},
};
</script>
