<template>
  <div class="file-time">
    <div class="file-time__title">
      <template v-if="progress === 100">
        Seeding
      </template>
      <template v-else>
        {{ [seconds, "seconds"] | duration("humanize", true) }}
      </template>
    </div>
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
      seconds: 0,
    };
  },
  computed: {
    ...mapState("downloadEvents", ["downloadEvent"]),
  },
  watch: {
    downloadEvent(newEvent) {
      if (this.file.FileHash === newEvent.FileHash) {
        this.seconds = this.file.FileSize / newEvent.Bandwidth;
        this.bandwidth = newEvent.Bandwidth;
        this.progress = newEvent.Progress * 100;
      }
    },
  },

  methods: {},
};
</script>
