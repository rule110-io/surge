<template>
  <div
    class="transfers-info-card"
    :class="{ 'transfers-info-card_active': active }"
  >
    <div class="transfers-info-card__header">
      <div class="transfers-info-card__title">{{ title }}</div>
      <div class="transfers-info-card__info"><slot name="info" /></div>
    </div>
    <div class="transfers-info-card__body"><slot name="body" /></div>
  </div>
</template>

<style lang="scss">
@import "./TransfersInfoCard";
</style>

<script>
export default {
  props: {
    active: {
      type: Boolean,
      default: false,
    },
    title: {
      type: String,
      default: "",
    },
  },
  data: () => {
    return {};
  },

  watch: {},
  mounted() {},
};
</script>
