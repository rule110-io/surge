<template>
  <div id="app" class="main">
    <div class="main__wrapper"><Sidebar /> <Dashboard /></div>
  </div>
</template>

<script>
import Dashboard from "./components/Dashboard/Dashboard";
import Sidebar from "./components/Sidebar/Sidebar";

import "./assets/scss/_main.scss";

export default {
  name: "app",
  components: {
    Dashboard,
    Sidebar,
  },
};
</script>
