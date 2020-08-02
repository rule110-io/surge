<template>
  <div class="dashboard">
    <Header />
    <router-view
      v-if="remoteFiles !== false && localFiles !== false"
    ></router-view>
    <Preloader v-else />
    <NetworkStats />
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

export default {
  components: { Header, Preloader, NetworkStats },
  data: () => {
    return {};
  },
  computed: {
    ...mapState("files", ["remoteFiles", "localFiles"]),
  },
  methods: {},
};
</script>
