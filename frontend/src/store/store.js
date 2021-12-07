import Vue from "vue";
import Vuex from "vuex";

import Notifications from "./modules/notifications";
import GlobalBandwidth from "./modules/globalBandwidth";
import Snackbar from "./modules/snackbar";
import DarkTheme from "./modules/darkTheme";
import Version from "./modules/version";
import PubKey from "./modules/pubKey";
import Tour from "./modules/tour";
import Topics from "./modules/topics";

import Files from "./modules/files";

Vue.use(Vuex);

export const store = new Vuex.Store({
  strict: false,
  modules: {
    notifications: Notifications,
    globalBandwidth: GlobalBandwidth,
    files: Files,
    snackbar: Snackbar,
    darkTheme: DarkTheme,
    version: Version,
    tour: Tour,
    pubKey: PubKey,
    topics: Topics,
  },
});

export default store;
