import "core-js/stable";
import "regenerator-runtime/runtime";
import Vue from "vue";
import App from "./App.vue";
import router from "./router/index.js";
import VueFeather from "vue-feather";
import VueMoment from "vue-moment";
import vueFilterPrettyBytes from "vue-filter-pretty-bytes";
import VueLodash from "vue-lodash";
import lodash from "lodash";
import VueTour from "vue-tour";
import VTooltip from "v-tooltip";
import VueClipboard from "vue-clipboard2";
import VueBus from "vue-bus";

import { store } from "./store/store.js";

Vue.config.productionTip = false;
Vue.config.devtools = true;

Vue.use(VueFeather);
Vue.use(vueFilterPrettyBytes);
Vue.use(VueMoment);
Vue.use(VueLodash, { lodash: lodash });
Vue.use(VueTour);
Vue.use(VTooltip);
Vue.use(VueClipboard);
Vue.use(VueBus);

new Vue({
  router,
  store,
  render: (h) => h(App),
  mounted() {
    this.$router.replace("/").catch(() => {});
  },
}).$mount("#app");
