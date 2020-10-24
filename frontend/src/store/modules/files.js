const getDefaultState = () => {
  return {
    localFiles: false,
    remoteFiles: false,
    localCount: 0,
    remoteCount: 0,
    localPages: 0,
    remotePages: 0,
    localFilesConfig: {
      skip: 0,
      get: 8,
    },
    remoteFilesConfig: {
      search: "",
      skip: 0,
      get: 8,
    },
  };
};

const state = getDefaultState();

const mutations = {
  setLocalFiles(state, payload) {
    const { Result, Count } = payload;
    state.localFiles = Result;
    state.localCount = Count;
    state.localPages = Math.ceil(Count / state.localFilesConfig.get);
  },
  setRemoteFiles(state, payload) {
    const { Result, Count } = payload;
    state.remoteFiles = Result;
    state.remoteCount = Count;
    state.remotePages = Math.ceil(Count / state.remoteFilesConfig.get);
  },
  setRemoteFilesConfig(state, payload) {
    state.remoteFilesConfig = payload;
  },
  setLocalFilesConfig(state, payload) {
    state.localFilesConfig = payload;
  },
};

const actions = {
  fetchLocalFiles({ commit, state }) {
    const { skip, get } = state.localFilesConfig;

    window.backend.getLocalFiles(skip, get).then(({ Result, Count }) => {
      commit("setLocalFiles", { Result, Count });
    });
  },
  fetchRemoteFiles({ commit, state }) {
    const { search, skip, get } = state.remoteFilesConfig;

    window.backend
      .getRemoteFiles(search, skip, get)
      .then(({ Result, Count }) => {
        commit("setRemoteFiles", { Result, Count });
      });
  },
};

export default {
  namespaced: true,
  state,
  getters: {},
  actions,
  mutations,
};
