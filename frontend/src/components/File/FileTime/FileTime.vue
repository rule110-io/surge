<template>
  <div class="file-time">
    <div class="file-time__title">11 min</div>
    <div class="file-time__percent">80%</div>
  </div>
</template>

<style lang="scss">
@import "./FileTime.scss";
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
    return {};
  },
  methods: {},
};
</script>
