<template>
  <div class="navigation">
    <router-link
      :to="`/${item.route}`"
      :id="item.route"
      class="navigation__item"
      v-for="item in items"
      :key="item.route"
    >
      <Icon class="navigation__icon" :icon="item.icon"></Icon>
      {{ item.text }}
    </router-link>
  </div>
</template>

<style lang="scss">
@import "./Navigation.scss";
</style>

<script>
import Icon from "@/components/Icon/Icon";

export default {
  components: { Icon },
  data: () => {
    return {
      items: [
        {
          text: "Transfers",
          route: "transfers",
          icon: "TransfersIcon",
        },
        {
          text: "Discover",
          route: "discover",
          icon: "DiscoverIcon",
        },
      ],
    };
  },
  computed: {},

  methods: {},
};
</script>
