<template>
  <div class="file-size  text_wrap_none">
    {{ file.FileSize | prettyBytes(1) }}
  </div>
</template>

<style lang="scss">
@import "./FileSize.scss";
</style>

<script>
export default {
  components: {},
  props: {
    file: {
      type: Object,
      default: () => {},
    },
  },
  computed: {},
  data: () => {
    return {};
  },
  methods: {},
};
</script>
