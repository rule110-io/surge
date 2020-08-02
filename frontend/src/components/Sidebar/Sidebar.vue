<template>
  <div class="sidebar">
    <Logo class="sidebar__logo" />
    <div class="sidebar__nav">
      <router-link to="/search" class="sidebar__item">
        <feather class="sidebar__item-icon" type="search"></feather
      ></router-link>
      <router-link to="/download" class="sidebar__item">
        <feather class="sidebar__item-icon" type="play"></feather
      ></router-link>
      <router-link to="/" class="sidebar__item">
        <feather class="sidebar__item-icon" type="file"></feather>
      </router-link>
      <router-link to="/" class="sidebar__item">
        <feather class="sidebar__item-icon" type="trash"></feather
      ></router-link>
    </div>
    <router-link to="/" class="sidebar__item">
      <feather class="sidebar__item-icon" type="settings"></feather
    ></router-link>
  </div>
</template>

<style lang="scss">
@import "./Sidebar.scss";
</style>

<script>
import Logo from "@/assets/images/Logo.svg";

export default {
  components: { Logo },
};
</script>
