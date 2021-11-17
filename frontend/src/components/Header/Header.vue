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
      <Icon class="header__icon" icon="SettingsIcon" />
    </div>

    <Modal :show.sync="showAddFileModal">
      <template slot="title"> Add New File </template>
      <template slot="body">
        <ControlWrapper title="Topic name">
          <Select
            v-model="topicName"
            :items="topics"
            placeholder="Select topic"
          />
        </ControlWrapper>
      </template>
      <template slot="footer">
        <Button theme="text" size="md" @click="closeAddFileModal">Close</Button>
        <Button theme="default" size="md" @click="addFile(topicName)"
          >Upload File</Button
        >
      </template>
    </Modal>
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
import Modal from "@/components/Modals/Modal/Modal";
import Select from "@/components/Controls/Select/Select";
import ControlWrapper from "@/components/Controls/ControlWrapper/ControlWrapper";

import Logo from "@/assets/icons/Logo.svg";

export default {
  components: {
    Logo,
    Navigation,
    Button,
    Divider,
    Icon,
    Notifications,
    Modal,
    ControlWrapper,
    Select,
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
    ...mapState("darkTheme", ["darkTheme"]),
    ...mapState("topics", ["topics"]),
    currentRoute() {
      return this.$route.name;
    },
  },
  methods: {
    addFile(topicName) {
      this.seedFile(topicName);
      this.closeAddFileModal();
      this.clearAddFileModal();
    },
    openAddFileModal() {
      this.showAddFileModal = true;
    },
    closeAddFileModal() {
      this.showAddFileModal = false;
    },
    clearAddFileModal() {
      this.topicName = "";
    },
    seedFile(topicName) {
      window.go.surge.MiddlewareFunctions.SeedFile(topicName).then(() => {
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
  },
};
</script>
