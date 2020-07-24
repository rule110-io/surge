<template>
  <div class="select__wrapper" :class="`select__wrapper_${type}`">
    <div
      class="select__button"
      :class="[open ? 'select__button_active' : null, `select__button_${type}`]"
      @click="toggleSelect()"
    >
      {{ $t(activeItem) }}
      <span
        class="select__icon fe fe-chevron-down"
        :class="open ? 'select__icon_open' : null"
      />
    </div>

    <ul class="select__list" :class="[open ? 'select__list_open' : null, `select__list_${type}`]">
      <li
        v-for="item in items"
        :key="item"
        class="select__item"
        :class="`select__item_${type}`"
        @click="setSelect(item),toggleSelect()"
      >
        {{ $t(item) }}
      </li>
    </ul>
  </div>
</template>

<style lang="scss">
@import "./Select.scss";
</style>

<script>
export default {
  props: {
    items: {
      type: Array,
      default: () => []
    },
    activeItem: {
      type: String,
      default: ''
    },
    type: {
      type: String,
      default: ''
    }
  },
  data: () => {
    return {
      open: false
    }
  },
  computed: {
  },
  methods: {
    toggleSelect () {
      this.open = !this.open
    },
    setSelect (item) {
      this.$emit('update', item)
    }
  }
}
</script>
