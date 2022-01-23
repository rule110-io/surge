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
  data: () => {
    return {};
  },
  computed: {},
  mounted() {},
  methods: {},
};
</script>
