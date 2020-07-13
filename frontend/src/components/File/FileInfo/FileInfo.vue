<template>
  <div class="file-info">
    <feather class="file-info__icon" type="file"></feather>
    <div class="file-info__right">
      <div class="file-info__title text_wrap_none">{{ file.FileName }}</div>
      <div class="file-info__size">{{ file.FileSize | prettyBytes(1) }}</div>
    </div>
  </div>
</template>

<style lang="scss">
@import "./FileInfo.scss";
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
