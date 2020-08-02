<template>
  <div class="preloader">
    <Logo class="preloader__logo" />
  </div>
</template>

<style lang="scss">
@import "./Preloader.scss";
</style>

<script>
import Logo from "@/assets/images/Logo.svg";

export default {
  components: { Logo },
  data: () => {
    return {};
  },
  methods: {},
};
</script>
