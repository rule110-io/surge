<template>
  <div class="dashboard">
    <Header />
    <router-view
      v-if="remoteFiles !== false && localFiles !== false"
    ></router-view>
  </div>
</template>

<style lang="scss">
@import "./Dashboard.scss";
</style>

<script>
import { mapState } from "vuex";

import Header from "@/components/Header/Header";

export default {
  components: { Header },
  data: () => {
    return {};
  },
  computed: {
    ...mapState("files", ["remoteFiles", "localFiles"]),
  },
  methods: {},
};
</script>
