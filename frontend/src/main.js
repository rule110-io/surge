import "core-js/stable";
import "regenerator-runtime/runtime";
import Vue from "vue";
import App from "./App.vue";
import VueRouter from "vue-router";
import VueFeather from "vue-feather";

import * as Wails from "@wailsapp/runtime";

Vue.config.productionTip = false;
Vue.config.devtools = true;

Vue.use(VueRouter);
Vue.use(VueFeather);

const Foo = { template: "<div>foo</div>" };
const Bar = { template: "<div>bar</div>" };

const routes = [
  { path: "/foo", component: Foo },
  { path: "/bar", component: Bar },
];

const router = new VueRouter({
  routes,
});

Wails.Init(() => {
  new Vue({
    router,
    render: (h) => h(App),
  }).$mount("#app");
});
