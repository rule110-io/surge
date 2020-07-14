<template>
  <div class="file-status">
    <div class="file-status__speed">{{ bandwith | prettyBytes(1) }}/s</div>
    <div class="file-status__progress">
      <div
        class="file-status__progress-current"
        :style="{ width: `${progress}%` }"
      ></div>
      <div class="file-status__progress-default"></div>
    </div>
  </div>
</template>

<style lang="scss">
@import "./FileStatus.scss";
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
  data() {
    return {
      progress: 0,
      bandwith: 0,
    };
  },
  computed: {
    ...mapState("downloadEvents", ["downloadEvent"]),
  },
  watch: {
    downloadEvent(newEvent) {
      if (this.file.FileHash === newEvent.FileHash) {
        this.bandwidth = newEvent.Bandwith;
        this.progress = newEvent.Progress * 100;
      }
    },
  },
  mounted() {},
  methods: {},
};
</script>
