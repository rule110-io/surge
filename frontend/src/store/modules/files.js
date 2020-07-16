const getDefaultState = () => {
  return {
    localFiles: false,
    remoteFiles: false,
  };
};

const state = getDefaultState();

const mutations = {
  setLocalFiles(state, localFiles) {
    state.localFiles = localFiles;
    console.log(localFiles);
  },
  setRemoteFiles(state, remoteFiles) {
    state.remoteFiles = remoteFiles;
    console.log(remoteFiles);
  },
};

const actions = {
  initLocalFiles({ commit }) {
    window.backend.getLocalFiles().then((result) => {
      commit("setLocalFiles", result);
    });
  },
  initRemoteFiles({ commit }) {
    window.backend.getRemoteFiles().then((result) => {
      commit("setRemoteFiles", result);
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
