<template>
  <div class="control-wrapper" :class="`col_${col}`">
    <label v-if="title" class="control-wrapper__title">{{ title }}</label>
    <div class="control-wrapper__descr" v-if="$slots.descr">
      <slot name="descr" />
    </div>
    <div v-if="type === 'row'" class="control-wrapper__row"><slot /></div>
    <slot v-else />
    <div
      v-if="showError"
      class="control-wrapper__error"
      :class="showError ? 'control-wrapper__error_visible' : null"
    >
      {{ errorMsg }}
    </div>
  </div>
</template>

<style lang="scss">
@import "./ControlWrapper.scss";
</style>

<script>
export default {
  props: {
    title: {
      type: String,
      default: "",
    },
    col: {
      type: String,
      default: "1",
    },
    type: {
      type: String,
      default: "column",
    },
    errorMsg: {
      type: String,
      default: "",
    },
    showError: {
      type: Boolean,
      default: false,
    },
  },
  methods: {},
};
</script>
