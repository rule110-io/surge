import "core-js/stable";
import "regenerator-runtime/runtime";
import Vue from "vue";
import App from "./App.vue";
import router from "./router/index";
import VueFeather from "vue-feather";
import vueFilterPrettyBytes from "vue-filter-pretty-bytes";

import * as Wails from "@wailsapp/runtime";

Vue.config.productionTip = false;
Vue.config.devtools = true;

Vue.use(VueFeather);
Vue.use(vueFilterPrettyBytes);

Wails.Init(() => {
  new Vue({
    router,
    render: (h) => h(App),
  }).$mount("#app");
});
