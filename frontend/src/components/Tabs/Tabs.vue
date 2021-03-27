<template>
  <div class="tabs">
    <div v-if="title" class="tabs__title">{{ title }}</div>
    <div class="tabs__items">
      <div
        :class="tab === activeTab ? `tabs__item_active` : null"
        class="tabs__item"
        v-for="(tab, i) in tabs"
        :key="i"
        @click="setActiveTab(tab)"
      >
        {{ tab }}
      </div>
    </div>
  </div>
</template>

<style lang="scss">
@import "./Tabs";
</style>

<script>
export default {
  props: {
    title: {
      type: String,
      default: "",
    },
    activeTab: {
      type: String,
      default: "",
    },
    tabs: {
      type: Array,
      default: () => [],
    },
  },
  components: {},
  data: () => {
    return {};
  },
  methods: {
    setActiveTab(tab) {
      this.$emit("update:activeTab", tab);
    },
  },
};
</script>
