<template>
  <div class="file-up text_wrap_none">
    {{ uploadBandwidth | prettyBytes(0) }}/s
  </div>
</template>

<style lang="scss">
@import "./FileUp.scss";
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
      uploadBandwidth: 0,
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
        this.uploadBandwidth = newEvent.UploadBandwidth;
      }
    },
  },

  methods: {},
};
</script>
