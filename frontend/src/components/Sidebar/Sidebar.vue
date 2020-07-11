<template>
  <div class="sidebar">
    <div class="sidebar__logo" />
    <div class="sidebar__nav">
      <router-link to="/download" class="sidebar__item">
        <feather class="sidebar__item-icon" type="play"></feather
      ></router-link>
      <router-link to="/folder" class="sidebar__item">
        <feather class="sidebar__item-icon" type="file"></feather>
      </router-link>
      <router-link to="/trash" class="sidebar__item">
        <feather class="sidebar__item-icon" type="trash"></feather
      ></router-link>
    </div>
    <router-link to="/settings" class="sidebar__item">
      <feather class="sidebar__item-icon" type="settings"></feather
    ></router-link>
  </div>
</template>

<style lang="scss">
@import "./Sidebar.scss";
</style>

<script>
export default {
  data() {
    return {};
  },
  methods: {},
};
</script>
