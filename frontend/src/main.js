import "core-js/stable";
import "regenerator-runtime/runtime";
import Vue from "vue";
import App from "./App.vue";
import router from "./router/index";
import VueFeather from "vue-feather";

import * as Wails from "@wailsapp/runtime";

Vue.config.productionTip = false;
Vue.config.devtools = true;

Vue.use(VueFeather);

Wails.Init(() => {
  new Vue({
    router,
    render: (h) => h(App),
    mounted() {
      this.$router.replace("/");
    },
  }).$mount("#app");
});
