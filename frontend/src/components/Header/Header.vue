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
      <template slot="title">
        Add New File
      </template>
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
import Input from "@/components/Controls/Input/Input";
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
    Input,
    ControlWrapper,
    Select,
  },
  data: () => {
    return {
      topicName: null,
      showAddFileModal: false,
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
    ...mapState("topics", ["topics"]),
    currentRoute() {
      return this.$route.name;
    },
  },
  created() {
    this.initRemoteSearch();
    this.initLocalSearch();
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
  },
};
</script>
