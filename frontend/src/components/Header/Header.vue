<template>
  <div class="header">
    <div class="header__left">
      <div :class="['header__search', focus ? 'header__search_active' : '']">
        <input
          type="text"
          class="header__search-input"
          placeholder="Search for remote files..."
          @focus="focus = true"
          @blur="focus = false"
          v-model.trim="searchQuery"
          @input="search(searchQuery)"
        />
        <div class="header__search-right">
          <feather class="header__search-icon" type="search"></feather>
        </div>
      </div>
    </div>
    <div class="header__right">
      <div class="header__item" @click="toggleTheme">
        <feather
          class="header__item-icon"
          v-if="!darkTheme"
          type="moon"
        ></feather>
        <feather class="header__item-icon" v-else type="sun"></feather>
      </div>
      <router-link to="/settings" class="header__item">
        <feather class="header__item-icon" type="settings"></feather
      ></router-link>
      <div
        class="header__item"
        @click="toggleNotifications"
        v-on-clickaway="closeNotifications"
      >
        <span
          :class="['header__badge', counter > 0 ? 'header__badge_visible' : '']"
          >{{ counter }}</span
        >
        <feather
          :class="[
            'header__item-icon',
            open > 0 ? 'header__item-icon_active' : '',
          ]"
          type="bell"
        ></feather>
        <Notifications @click.native.stop.prevent />
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
  data: () => {
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
    ...mapState("darkTheme", ["darkTheme"]),
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
    toggleTheme() {
      this.$store.dispatch("darkTheme/toggleDarkTheme");
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
