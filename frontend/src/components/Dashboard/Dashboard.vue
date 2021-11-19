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
  </div>
</template>

<style lang="scss">
@import "./Dashboard.scss";
</style>

<script>
import { mapState } from "vuex";

import Header from "@/components/Header/Header";
import NetworkStats from "@/components/NetworkStats/NetworkStats";
import AddTopicModal from "@/components/Modals/AddTopicModal/AddTopicModal";
import AddFileModal from "@/components/Modals/AddFileModal/AddFileModal";

export default {
  components: { Header, NetworkStats, AddTopicModal, AddFileModal },
  data: () => {
    return {};
  },
  computed: {
    ...mapState("files", ["remoteFiles", "localFiles"]),
  },
  methods: {},
};
</script>
