import axios from "axios";

const compareVersions = require("compare-versions");

const state = () => ({
  currentVersion: require("../../../package.json").version,
  remoteVersion: false,
  isNewVersion: false,
});

const mutations = {
  setRemoteVersion(state, version) {
    state.remoteVersion = version;
  },
  setIsNewVersion(state, bool) {
    state.isNewVersion = bool;
  },
};

const getters = {
  getCurrentVersion(state) {
    return state.currentVersion;
  },
  getRemoteVersion(state) {
    return state.remoteVersion;
  },
  getIsNewVersion(state) {
    return state.isNewVersion;
  },
};

const actions = {
  async updateRemoteVersion({ commit }) {
    const data = await axios(
      "https://api.github.com/repos/rule110-io/surge/releases"
    );

    const currentVersion = this.state.version.currentVersion;

    const releases = data.data;

    const remoteVersion = releases ? releases[0].tag_name : currentVersion;
    const isNewVersion = compareVersions.compare(
      remoteVersion,
      currentVersion,
      ">"
    );

    commit("setRemoteVersion", remoteVersion);
    commit("setIsNewVersion", isNewVersion);
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
