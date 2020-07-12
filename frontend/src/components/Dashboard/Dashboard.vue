<template>
  <div class="dashboard">
    <Header />
    <router-view></router-view>
  </div>
</template>

<style lang="scss">
@import "./Dashboard.scss";
</style>

<script>
import Header from "@/components/Header/Header";

export default {
  components: { Header },
  data() {
    return {};
  },
  methods: {},
};
</script>
