<template>
  <component
    :is="type"
    class="button"
    :class="[
      `button_theme_${theme}`,
      active ? `button_theme_${theme}_active` : null,
      `button_size_${size}`,
    ]"
    :to="to"
    :disabled="disabled"
    v-on="!disabled ? $listeners : false"
  >
    <Icon :icon="icon" v-if="icon" class="button__icon"> </Icon>
    <slot></slot>
  </component>
</template>

<style lang="scss">
@import "./Button";
</style>

<script>
import Icon from "@/components/Icon/Icon";

export default {
  components: { Icon },
  props: {
    theme: {
      type: String,
      default: "default",
    },
    size: {
      type: String,
      default: "lg",
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
    active: {
      type: Boolean,
      default: false,
    },
  },
  data: () => {
    return {};
  },
  mounted() {},
  methods: {},
};
</script>
