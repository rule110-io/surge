import Vue from "vue";
import Vuex from "vuex";

import Notifications from "./modules/notifications";
import ClientStatus from "./modules/clientStatus";
import GlobalBandwidth from "./modules/globalBandwidth";
import Snackbar from "./modules/snackbar";
import DarkTheme from "./modules/darkTheme";
import Version from "./modules/version";
import PubKey from "./modules/pubKey";
import Tour from "./modules/tour";

import Files from "./modules/files";

Vue.use(Vuex);

export const store = new Vuex.Store({
  strict: false,
  modules: {
    notifications: Notifications,
    clientStatus: ClientStatus,
    globalBandwidth: GlobalBandwidth,
    files: Files,
    snackbar: Snackbar,
    darkTheme: DarkTheme,
    version: Version,
    tour: Tour,
    pubKey: PubKey,
  },
});

export default store;
