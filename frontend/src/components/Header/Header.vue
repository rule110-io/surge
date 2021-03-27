<template>
  <header class="header">
    <div class="header__left">
      <Logo class="header__logo" />
      <Navigation />
    </div>

    <div class="header__right">
      <Input
        class="header__search"
        :value="searchQuery"
        icon="SearchIcon"
        placeholder="Search or enter file hash"
        @update="
          currentRoute === 'download'
            ? localSearch(searchQuery)
            : remoteSearch(searchQuery)
        "
      />
      <Button theme="primary" class="header__button">Add file</Button>
      <Divider />
    </div>
  </header>
</template>

<style lang="scss">
@import "./Header.scss";
</style>

<script>
import { mapState } from "vuex";
import { mixin as clickaway } from "vue-clickaway";

import Navigation from "@/components/Navigation/Navigation";
import Input from "@/components/Controls/Input/Input";
import Button from "@/components/Button/Button";
import Divider from "@/components/Divider/Divider";

import Logo from "@/assets/icons/Logo.svg";

export default {
  components: { Logo, Navigation, Input, Button, Divider },
  mixins: [clickaway],
  data: () => {
    return {
      active: true,
      focus: false,
      searchQuery: "",
      remoteSearch: () => {},
      localSearch: () => {},
    };
  },
  computed: {
    ...mapState("notifications", ["counter", "open"]),
    ...mapState("files", ["remoteFilesConfig", "localFilesConfig"]),
    ...mapState("darkTheme", ["darkTheme"]),
    currentRoute() {
      return this.$route.name;
    },
  },
  created() {
    this.initRemoteSearch();
    this.initLocalSearch();
  },
  methods: {
    initRemoteSearch() {
      this.remoteSearch = this._.debounce((search) => {
        if (this.currentRoute !== "search") {
          this.$router.replace("/search");
        }

        let newConfig = Object.assign({}, this.remoteFilesConfig);
        newConfig.skip = 0;
        newConfig.search = search;

        this.$store.commit("files/setRemoteFilesConfig", newConfig);
        this.$store.dispatch("files/fetchRemoteFiles");
      }, 500);
    },
    initLocalSearch() {
      this.localSearch = this._.debounce((search) => {
        let newConfig = Object.assign({}, this.localFilesConfig);
        newConfig.skip = 0;
        newConfig.search = search;

        this.$store.commit("files/setLocalFilesConfig", newConfig);
        this.$store.dispatch("files/fetchLocalFiles");
      }, 500);
    },
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
