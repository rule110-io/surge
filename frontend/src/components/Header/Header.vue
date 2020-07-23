<template>
  <div class="header">
    <div class="header__left">
      <div class="header__search" :class="focus ? 'header__search_active' : ''">
        <feather class="header__search-icon" type="search"></feather>
        <input
          type="text"
          class="header__search-input"
          placeholder="Search for remote files..."
          @focus="focus = true"
          @blur="focus = false"
          v-model.trim="searchQuery"
          @input="search(searchQuery)"
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
      <div class="header__file" @click="seedFile">
        <feather class="header__file-icon" type="plus"></feather>
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
      searchQuery: "",
      search: () => {},
    };
  },
  computed: {
    ...mapState("notifications", ["counter", "open"]),
    ...mapState("files", ["remoteFilesConfig"]),
  },
  created() {
    this.search = this._.debounce((search) => {
      if (this.$router.currentRoute.name !== "search") {
        this.$router.replace("/search");
      }

      let newConfig = Object.assign({}, this.remoteFilesConfig);
      newConfig.skip = 0;
      newConfig.search = search;

      this.$store.commit("files/setRemoteFilesConfig", newConfig);
      this.$store.dispatch("files/fetchRemoteFiles");
    }, 500);
  },
  mounted() {},
  methods: {
    toggleNotifications() {
      this.$store.commit("notifications/toggleNotifications", !this.open);
    },
    closeNotifications() {
      this.$store.commit("notifications/toggleNotifications", false);
    },
    seedFile() {
      window.backend.seedFile().then(() => {});
    },
  },
};
</script>
