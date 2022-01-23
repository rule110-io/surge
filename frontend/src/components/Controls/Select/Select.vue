<template>
  <div class="select__wrapper" v-on-clickaway="closeSelect">
    <div
      :class="['select__button', open ? 'select__button_active' : null]"
      @click="toggleSelect()"
    >
      <span
        v-if="!isPlaceholder"
        class="select__value"
        v-text="itemText ? value[itemText] : value"
      />
      <span class="select__placeholder" v-else>{{ placeholder }}</span>

      <Icon
        :class="['select__icon', open ? 'select__icon_open' : null]"
        icon="SelectIcon"
      />
    </div>

    <Dropdown :open.sync="open" theme="light">
      <ul class="select__list">
        <li
          v-for="item in items"
          :key="item"
          :class="`select__item`"
          @click="setSelect(item), toggleSelect()"
          v-text="itemText ? item[itemText] : item"
        ></li>
      </ul>
    </Dropdown>
  </div>
</template>

<style lang="scss">
@import "./Select.scss";
</style>

<script>
import { mixin as clickaway } from "vue-clickaway";

import Icon from "@/components/Icon/Icon";
import Dropdown from "@/components/Dropdown/Dropdown";

export default {
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    value: {
      type: String,
      default: "",
    },
    itemText: {
      type: String,
      default: "",
    },
    placeholder: {
      type: String,
      default: "Click to select",
    },
  },
  components: { Icon, Dropdown },
  mixins: [clickaway],
  data: () => {
    return {
      open: false,
    };
  },
  computed: {
    isPlaceholder() {
      const { value, itemText } = this;

      if (itemText.length) {
        const hasValue = this._.has(value, itemText);
        return hasValue
          ? this._.isEmpty(this._.toString(value[itemText]))
          : true;
      } else {
        return this._.isEmpty(this._.toString(value));
      }
    },
  },
  methods: {
    toggleSelect() {
      this.open = !this.open;
    },
    setSelect(item) {
      this.$emit("input", item);
    },
    closeSelect() {
      if (!this.open) return;

      this.open = false;
    },
  },
};
</script>
