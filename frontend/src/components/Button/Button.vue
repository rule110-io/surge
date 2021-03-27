<template>
  <component
    :is="type"
    class="button"
    :class="[`button_theme_${theme}`, `button_size_${size}`]"
    :to="to"
    :disabled="disabled"
    v-on="!disabled ? $listeners : false"
  >
    <component :is="icon" v-if="icon" class="button__icon"> </component>
    <slot></slot>
  </component>
</template>

<style lang="scss">
@import "./Button";
</style>

<script>
export default {
  components: {},
  props: {
    theme: {
      type: String,
      default: "default",
    },
    size: {
      type: String,
      default: "md",
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    type: {
      type: String,
      default: "button",
    },
    to: {
      type: String,
      default: "",
    },
    icon: {
      type: String,
      default: "",
    },
  },
  data: () => {
    return {};
  },
  mounted() {},
  methods: {},
};
</script>
