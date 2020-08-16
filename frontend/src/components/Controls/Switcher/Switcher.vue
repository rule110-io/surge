<template>
  <div class="switch">
    <input
      v-model="radioButtonValue"
      class="switch__control"
      type="checkbox"
      :name="name"
    />
    <label class="switch__label" />
    <div class="switch__container"></div>
  </div>
</template>

<style lang="scss">
@import "./Switcher.scss";
</style>

<script>
export default {
  props: {
    value: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: "",
    },
    label: {
      type: String,
      default: "",
    },
  },
  computed: {
    radioButtonValue: {
      get() {
        return this.value;
      },
      set() {
        this.$emit("change", this.label);
      },
    },
  },
};
</script>
