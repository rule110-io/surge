<template>
  <div class="file-up text_wrap_none">
    <template v-if="uploadBandwidth > 0">
      {{ uploadBandwidth | prettyBytes(0) }}/s
    </template>
    <template v-else> - </template>
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
    ...mapState("globalBandwidth", ["statusBundle"]),
  },
  watch: {
    statusBundle(newEvent) {
      if (!this.file) return;

      const { FileHash } = this.file;
      const newFileHash = this._.find(newEvent, { FileHash });
      const isNewFileHash = !this._.isEmpty(newFileHash);

      if (isNewFileHash) {
        this.uploadBandwidth = newFileHash.UploadBandwidth;
      }
    },
  },

  methods: {},
};
</script>
