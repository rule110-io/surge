<template>
  <div
    style="
      -moz-user-select: none;
      -webkit-user-select: none;
      -ms-user-select: none;
      user-select: none;
      -o-user-select: none;
    "
    id="app"
    class="main"
    :class="darkTheme ? 'theme_dark' : null"
  >
    <router-view></router-view>
  </div>
</template>

<script>
import "./assets/scss/_main.scss";

import { mapState } from "vuex";

export default {
  name: "app",
  components: {},
  computed: {
    ...mapState("darkTheme", ["darkTheme"]),
  },
};
</script>
