<template>
  <div class="header">
    <div class="header__left">
      <div class="header__search" :class="focus ? 'header__search_active' : ''">
        <feather class="header__search-icon" type="search"></feather>
        <input
          type="text"
          class="header__search-input"
          placeholder="Search for files and more..."
          @focus="focus = true"
          @blur="focus = false"
        />
      </div>
    </div>
    <div class="header__right">
      <router-link to="/settings" class="header__item">
        <feather class="header__item-icon" type="settings"></feather
      ></router-link>
      <router-link to="/notifications" class="header__item">
        <feather class="header__item-icon" type="bell"></feather
      ></router-link>
      <div class="header__avatar">
        <div
          class="header__status"
          :class="active ? 'header__status_active' : 'header__status_inactive'"
        ></div>
      </div>
    </div>
  </div>
</template>

<style lang="scss">
@import "./Header.scss";
</style>

<script>
export default {
  data() {
    return {
      active: true,
      focus: false,
    };
  },
  methods: {},
};
</script>
