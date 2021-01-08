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
      const { FileHash } = this.file;
      var self = this;
      newEvent.forEach(function (file) {
        if (FileHash === file.FileHash) {
          self.downloadBandwidth = file.DownloadBandwidth;
        }
      });
    },
  },

  methods: {},
};
</script>
