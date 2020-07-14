<template>
  <div class="file-status">
    <div class="file-status__speed">34.5 MB/s</div>
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
export default {
  props: {
    file: {
      type: Object,
      default: () => {},
    },
  },
  data() {
    return {
      progress: 70,
    };
  },
  methods: {},
};
</script>
