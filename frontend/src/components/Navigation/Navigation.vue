<template>
  <div class="navigation">
    <router-link
      :to="`/${item.route}`"
      :id="item.route"
      class="navigation__item"
      v-for="item in items"
      :key="item.route"
    >
      <component class="navigation__icon" :is="item.icon"></component>
      {{ item.text }}
    </router-link>
  </div>
</template>

<style lang="scss">
@import "./Navigation.scss";
</style>

<script>
import TransfersIcon from "@/assets/icons/TransfersIcon.svg";
import DiscoverIcon from "@/assets/icons/DiscoverIcon.svg";

export default {
  components: { TransfersIcon, DiscoverIcon },
  data: () => {
    return {
      items: [
        {
          text: "Transfers",
          route: "download",
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
