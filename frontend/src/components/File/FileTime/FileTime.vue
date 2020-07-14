<template>
  <div class="file-time">
    <div class="file-time__title">11 min</div>
    <div class="file-time__percent">{{ progress.toFixed(2) }}%</div>
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
  data() {
    return {
      progress: 0,
      bandwidth: 0,
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

  methods: {},
};
</script>
