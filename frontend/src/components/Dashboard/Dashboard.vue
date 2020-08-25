<template>
  <div class="dashboard">
    <Snackbar />
    <Header />
    <router-view
      v-if="remoteFiles !== false && localFiles !== false"
    ></router-view>
    <Preloader v-else />
    <NetworkStats />
    <StartDownloadModal />
  </div>
</template>

<style lang="scss">
@import "./Dashboard.scss";
</style>

<script>
import { mapState } from "vuex";

import Header from "@/components/Header/Header";
import Preloader from "@/components/Preloader/Preloader";
import NetworkStats from "@/components/NetworkStats/NetworkStats";
import Snackbar from "@/components/Snackbar/Snackbar.vue";
import StartDownloadModal from "@/components/Modals/StartDownloadModal/StartDownloadModal";

export default {
  components: { Header, Preloader, NetworkStats, Snackbar, StartDownloadModal },
  data: () => {
    return {};
  },
  computed: {
    ...mapState("files", ["remoteFiles", "localFiles"]),
  },
  methods: {},
};
</script>
