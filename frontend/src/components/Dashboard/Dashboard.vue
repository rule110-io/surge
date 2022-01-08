<template>
  <div class="dashboard">
    <Header />
    <div class="dashboard__content">
      <router-view
        v-if="remoteFiles !== false && localFiles !== false"
      ></router-view>
    </div>

    <NetworkStats />
    <AddTopicModal openModalEvent="openAddTopicModal" />
    <AddFileModal openModalEvent="openAddFileModal" />
    <RemoveFileModal openModalEvent="openRemoveFileModal" />
    <SettingsModal openModalEvent="openSettingsModal" />
  </div>
</template>

<style lang="scss">
@import "./Dashboard.scss";
</style>

<script>
import { mapState, mapMutations, mapActions } from "vuex";

import Header from "@/components/Header/Header";
import NetworkStats from "@/components/NetworkStats/NetworkStats";
import AddTopicModal from "@/components/Modals/AddTopicModal/AddTopicModal";
import AddFileModal from "@/components/Modals/AddFileModal/AddFileModal";
import RemoveFileModal from "@/components/Modals/RemoveFileModal/RemoveFileModal";
import SettingsModal from "@/components/Modals/SettingsModal/SettingsModal";

export default {
  components: {
    Header,
    NetworkStats,
    AddTopicModal,
    AddFileModal,
    RemoveFileModal,
    SettingsModal,
  },
  data: () => {
    return {};
  },
  computed: {
    ...mapState("files", ["remoteFiles", "localFiles", "selectedFiles"]),
  },
  watch: {
    selectedFiles(newFiles) {
      if (!newFiles.length) {
        this.closeDetails();
      }
    },
    $route() {
      this.closeDetails();
    },
  },
  methods: {
    ...mapMutations({
      setFileDetails: "files/setFileDetails",
      setFileSpeed: "files/setFileSpeed",
    }),
    ...mapActions({
      clearSelectedFiles: "files/clearSelectedFiles",
    }),
    closeDetails() {
      this.setFileDetails(false);
      this.setFileSpeed(false);
    },
  },
};
</script>
