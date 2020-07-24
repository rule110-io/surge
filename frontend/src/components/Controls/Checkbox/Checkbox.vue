<template>
  <div class="checkbox">
    <input
      v-model="radioButtonValue"
      class="checkbox__control"
      type="checkbox"
      :checked="value"
      :name="name"
    >
    <label class="checkbox__label" />
    <span class="checkbox__text">
      <slot />
    </span>
  </div>
</template>

<style lang="scss">
@import "./Checkbox.scss";
</style>

<script>
export default {
  props: {
    value: {
      type: Boolean,
      default: false
    },
    name: {
      type: String,
      default: ''
    },
    label: {
      type: String,
      default: ''
    }
  },
  computed: {
    radioButtonValue: {
      get () {
        return this.value
      },
      set () {
        this.$emit('change', this.label)
      }
    }
  }
}
</script>
