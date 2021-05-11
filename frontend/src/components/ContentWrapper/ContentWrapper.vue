<template>
  <div class="content-wrapper" :class="`content-wrapper_type_${type}`">
    <slot></slot>
  </div>
</template>

<style lang="scss">
@import "./ContentWrapper";
</style>

<script>
export default {
  props: {
    type: {
      type: String,
      default: "",
    },
  },
};
</script>
