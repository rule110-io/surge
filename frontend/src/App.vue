<template>
  <div id="app" class="main" :class="darkTheme ? 'theme_dark' : null">
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
