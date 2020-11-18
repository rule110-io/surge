<template>
  <div
    class="file-hash text_wrap_none"
    v-tooltip="{
      content: hash,
      placement: 'bottom-start',
      offset: 5,
    }"
  >
    {{ hash }}
  </div>
</template>

<style lang="scss">
@import "./FileHash.scss";
</style>

<script>
export default {
  components: {},
  props: {
    hash: {
      type: String,
      default: "",
    },
  },
};
</script>
