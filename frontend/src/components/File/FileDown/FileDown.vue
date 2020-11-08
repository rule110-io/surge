<template>
  <div class="file-down text_wrap_none">
    {{ downloadBandwidth | prettyBytes(0) }}/s
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
    ...mapState("downloadEvents", ["downloadEvent"]),
  },
  watch: {
    downloadEvent(newEvent) {
      const { FileHash } = this.file;
      if (FileHash === newEvent.FileHash) {
        this.downloadBandwidth = newEvent.DownloadBandwidth;
      }
    },
  },

  methods: {},
};
</script>
