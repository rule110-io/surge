<template>
  <div
    class="checkbox"
    :class="[divider ? 'checkbox_divider' : null, `checkbox_theme_${theme}`]"
  >
    <input
      :disabled="disabled"
      class="checkbox__control"
      type="checkbox"
      :checked="isChecked"
      :value="value"
      @change="updateInput"
    />
    <label class="checkbox__label" :class="`checkbox__label_theme_${theme}`" />
    <span v-if="title || $slots.title" class="checkbox__title">
      <slot v-if="$slots.title" name="title" />
      <template v-else>
        {{ title }}
      </template>
    </span>
  </div>
</template>

<style lang="scss">
@import "./Checkbox.scss";
</style>

<script>
export default {
  model: {
    prop: "modelValue",
    event: "change",
  },
  props: {
    value: { type: [String, Number, Boolean, Object], required: true },
    modelValue: { type: [String, Array, Boolean, Object], default: "" },
    title: { type: String, default: "" },
    theme: { type: String, default: "" },
    trueValue: { type: Boolean, default: true },
    disabled: { type: Boolean, default: false },
    falseValue: { type: Boolean, default: false },
    divider: { type: Boolean, default: false },
  },
  computed: {
    isChecked() {
      if (Array.isArray(this.modelValue)) {
        return (
          this._.filter(this.modelValue, (x) => this._.isEqual(x, this.value))
            .length > 0
        );
      }
      return this.modelValue === this.trueValue;
    },
  },
  methods: {
    updateInput(event) {
      const isChecked = event.target.checked;
      if (Array.isArray(this.modelValue)) {
        const newValue = [...this.modelValue];
        if (isChecked) {
          newValue.push(this.value);
        } else {
          newValue.splice(newValue.indexOf(this.value), 1);
        }
        this.$emit("change", newValue);
      } else {
        this.$emit("change", isChecked ? this.trueValue : this.falseValue);
      }
    },
  },
};
</script>
