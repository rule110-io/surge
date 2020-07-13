import Vue from "vue";
import Vuex from "vuex";
import Notifications from "./modules/notifications";

Vue.use(Vuex);

export const store = new Vuex.Store({
  modules: {
    notifications: Notifications,
  },
  getters: {
    runningOnWindows(state) {
      return state.OS.windows;
    },
    runningOnLinux(state) {
      return state.OS.linux;
    },
    runningOnMacOS(state) {
      return state.OS.macOS;
    },
  },
});

export default store;
