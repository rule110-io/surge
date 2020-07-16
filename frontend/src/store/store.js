import Vue from "vue";
import Vuex from "vuex";
import Notifications from "./modules/notifications";
import DownloadEvents from "./modules/downloadEvents";
import Files from "./modules/files";

Vue.use(Vuex);

export const store = new Vuex.Store({
  modules: {
    notifications: Notifications,
    downloadEvents: DownloadEvents,
    files: Files,
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
