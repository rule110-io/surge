<template>
  <div
    class="file-shared-badge"
    :class="
      ratio <= 1
        ? 'file-shared-badge_theme_default'
        : 'file-shared-badge_theme_good'
    "
    v-tooltip="{
      content: `This file shared ratio is ${ratio}x`,
      placement: 'bottom-center',
      offset: 5,
    }"
  >
    {{ ratio }}x
  </div>
</template>

<style lang="scss">
@import "./FileSharedBadge.scss";
</style>

<script>
export default {
  props: {
    ratio: {
      type: String,
      default: "0.00",
    },
  },
};
</script>
