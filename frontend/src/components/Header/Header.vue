<template>
  <header class="header">
    <div class="header__left">
      <Logo class="header__logo" />
      <Navigation />
    </div>

    <div class="header__right">
      <CustomInput
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
      <input placeholder="topic name" type="text" v-model="topicName" />
      <Button
        theme="primary"
        class="header__button"
        @click="seedFile(topicName)"
        >Add file</Button
      >
      <Divider />
      <Icon class="header__icon" icon="NotificationsIcon" />
      <Icon class="header__icon" icon="SettingsIcon" />
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
import CustomInput from "@/components/Controls/Input/Input";
import Button from "@/components/Button/Button";
import Divider from "@/components/Divider/Divider";
import Icon from "@/components/Icon/Icon";

import Logo from "@/assets/icons/Logo.svg";

export default {
  components: { Logo, Navigation, CustomInput, Button, Divider, Icon },
  mixins: [clickaway],
  data: () => {
    return {
      topicName: "",
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
    seedFile(topicName) {
      window.backend.MiddlewareFunctions.SeedFile(topicName).then(() => {
        this.$store.dispatch("files/fetchLocalFiles");
        this.$store.dispatch("files/fetchRemoteFiles");
      });
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
