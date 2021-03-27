<template>
  <div class="input" :class="focus ? 'input_focus' : null">
    <div v-if="icon" class="input__icon">
      <component :is="icon" />
    </div>
    <input
      class="input__controller"
      :class="[icon ? 'input__controller_icon' : null]"
      :type="type"
      :placeholder="placeholder"
      :value="value"
      @input="updateValue($event.target.value)"
      @focus="focus = true"
      @blur="focus = false"
    />
  </div>
</template>

<style lang="scss">
@import "./Input.scss";
</style>

<script>
import SearchIcon from "@/assets/icons/SearchIcon.svg";

export default {
  components: { SearchIcon },
  props: {
    type: {
      type: String,
      default: "text",
    },
    placeholder: {
      type: String,
      default: "Placeholder",
    },
    value: {
      type: [String, Number],
      default: "",
    },
    icon: {
      type: String,
      default: "",
    },
  },
  data: () => {
    return {
      focus: false,
    };
  },
  methods: {
    updateValue(value) {
      this.$emit("input", value);

      if (this.$listeners.update) {
        this.$listeners.update();
      }
    },
  },
};
</script>
