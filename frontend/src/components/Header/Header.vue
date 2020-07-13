<template>
  <div class="header">
    <div class="header__left">
      <div class="header__search" :class="focus ? 'header__search_active' : ''">
        <feather class="header__search-icon" type="search"></feather>
        <input
          type="text"
          class="header__search-input"
          placeholder="Search for files and more..."
          @focus="focus = true"
          @blur="focus = false"
        />
      </div>
    </div>
    <div class="header__right">
      <router-link to="/settings" class="header__item">
        <feather class="header__item-icon" type="settings"></feather
      ></router-link>
      <div
        class="header__item"
        @click="toggleNotifications"
        v-on-clickaway="closeNotifications"
      >
        <span
          class="header__badge"
          :class="counter > 0 ? 'header__badge_visible' : ''"
          >{{ counter }}</span
        >
        <feather
          class="header__item-icon"
          :class="open > 0 ? 'header__item-icon_active' : ''"
          type="bell"
        ></feather>
        <Notifications @click.native.stop.prevent />
      </div>
      <div class="header__avatar">
        <div
          class="header__status"
          :class="active ? 'header__status_active' : 'header__status_inactive'"
        ></div>
      </div>
    </div>
  </div>
</template>

<style lang="scss">
@import "./Header.scss";
</style>

<script>
import { mapState } from "vuex";
import { mixin as clickaway } from "vue-clickaway";

import Notifications from "@/components/Notifications/Notifications";

export default {
  components: { Notifications },
  mixins: [clickaway],
  data() {
    return {
      active: true,
      focus: false,
    };
  },
  computed: {
    ...mapState("notifications", ["counter", "open"]),
  },
  mounted() {
    this.enableNotifications();
  },
  methods: {
    enableNotifications() {
      window.wails.Events.On("notificationEvent", (title, text) => {
        const notification = { title, text };
        this.$store.commit("notifications/addNotification", notification);
      });
    },
    toggleNotifications() {
      this.$store.commit("notifications/toggleNotifications", !this.open);
    },
    closeNotifications() {
      this.$store.commit("notifications/toggleNotifications", false);
    },
  },
};
</script>
