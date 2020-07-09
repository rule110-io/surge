<template>
  <div id="app" class="main">
    <div class="main__wrapper"><HelloWorld /></div>
  </div>
</template>

<script>
import HelloWorld from "./components/HelloWorld.vue";
import "./assets/scss/_main.scss";

export default {
  name: "app",
  components: {
    HelloWorld,
  },
};
</script>
