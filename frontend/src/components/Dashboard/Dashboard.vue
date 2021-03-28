<template>
  <div class="dashboard">
    <Header />
    <div class="dashboard__content">
      <router-view
        v-if="remoteFiles !== false && localFiles !== false"
      ></router-view>
    </div>

    <NetworkStats />
  </div>
</template>

<style lang="scss">
@import "./Dashboard.scss";
</style>

<script>
import { mapState } from "vuex";

import Header from "@/components/Header/Header";
import NetworkStats from "@/components/NetworkStats/NetworkStats";

export default {
  components: { Header, NetworkStats },
  data: () => {
    return {};
  },
  computed: {
    ...mapState("files", ["remoteFiles", "localFiles"]),
  },
  methods: {},
};
</script>
