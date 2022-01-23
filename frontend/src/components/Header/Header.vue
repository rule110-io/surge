<template>
  <header class="header">
    <div class="header__left">
      <Logo class="header__logo" />
      <Navigation />
    </div>

    <div class="header__right">
      <Button theme="primary" class="header__button" @click="openAddFileModal"
        >Add file</Button
      >
      <Divider />
      <Icon
        class="header__icon"
        icon="NotificationsIcon"
        @click.native.stop="toggleNotifications"
      />
      <div class="header__notifications-marker" v-show="counter"></div>
      <Notifications />
      <Icon
        @click.native="openSettingsModal"
        class="header__icon"
        icon="SettingsIcon"
      />
    </div>
  </header>
</template>

<style lang="scss">
@import "./Header.scss";
</style>

<script>
import { mapState } from "vuex";

import Navigation from "@/components/Navigation/Navigation";
import Notifications from "@/components/Notifications/Notifications";
import Button from "@/components/Button/Button";
import Divider from "@/components/Divider/Divider";
import Icon from "@/components/Icon/Icon";

import Logo from "@/assets/icons/Logo.svg";

export default {
  components: {
    Logo,
    Navigation,
    Button,
    Divider,
    Icon,
    Notifications,
  },
  data: () => {
    return {
      topicName: null,
      showAddFileModal: false,
      active: true,
      focus: false,
    };
  },
  computed: {
    ...mapState("notifications", ["counter", "open"]),
    ...mapState("files", ["remoteFilesConfig", "localFilesConfig"]),
    currentRoute() {
      return this.$route.name;
    },
  },
  methods: {
    openAddFileModal() {
      this.$bus.$emit("openAddFileModal");
    },
    openSettingsModal() {
      this.$bus.$emit("openSettingsModal");
    },
    toggleNotifications() {
      this.$store.commit("notifications/toggleNotifications", !this.open);
    },
  },
};
</script>
