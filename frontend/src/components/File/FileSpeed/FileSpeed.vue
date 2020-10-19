<template>
  <div class="file-speed">
    <div class="file-speed__source">1 Source connected</div>
    <div>
      <div class="file-speed__item">
        Down: {{ downloadBandwidth | prettyBytes(1) }}/s
      </div>
      <div class="file-speed__item">
        Up: {{ uploadBandwidth | prettyBytes(1) }}/s
      </div>
    </div>
  </div>
</template>

<style lang="scss">
@import "./FileSpeed.scss";
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
      }
    },
  },
  mounted() {},
  methods: {},
};
</script>
