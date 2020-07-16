<template>
  <div class="dashboard">
    <Header />
    <router-view
      v-if="remoteFiles !== false && localFiles !== false"
    ></router-view>
    <Preloader v-else />
  </div>
</template>

<style lang="scss">
@import "./Dashboard.scss";
</style>

<script>
import { mapState } from "vuex";

import Header from "@/components/Header/Header";
import Preloader from "@/components/Preloader/Preloader";

export default {
  components: { Header, Preloader },
  data() {
    return {};
  },
  computed: {
    ...mapState("files", ["remoteFiles", "localFiles"]),
  },
  methods: {},
};
</script>
