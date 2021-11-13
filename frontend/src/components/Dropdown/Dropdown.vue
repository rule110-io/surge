<template>
  <div
    :class="[
      'dropdown',
      open ? 'dropdown_visible' : '',
      `dropdown_theme_${theme}`,
    ]"
  >
    <div class="dropdown__wrapper">
      <slot />
    </div>
  </div>
</template>

<style lang="scss">
@import "./Dropdown.scss";
</style>

<script>
import { mixin as clickaway } from "vue-clickaway";

export default {
  components: {},
  props: {
    open: {
      type: Boolean,
      required: true,
      default: false,
    },
    theme: {
      type: String,
      default: "dark",
    },
  },
  mixins: [clickaway],
  data: () => {
    return {};
  },
  computed: {},
  mounted() {},
  methods: {},
};
</script>
