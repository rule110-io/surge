import Vue from "vue";
import VueRouter from "vue-router";
import routes from "./routes";
Vue.use(VueRouter);

// configure router
const router = new VueRouter({
  routes,
  mode: "history",
  linkActiveClass: "active",
});

export default router;
