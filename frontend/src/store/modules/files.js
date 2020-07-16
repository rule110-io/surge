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
  },
  setRemoteFiles(state, remoteFiles) {
    state.remoteFiles = remoteFiles;
  },
};

const actions = {
  fetchLocalFiles({ commit }) {
    window.backend.getLocalFiles(0, 5).then(({ Result, Count }) => {
      commit("setLocalFiles", Result);
      console.log(Result, Count);
    });
  },
  fetchRemoteFiles({ commit }, payload) {
    let search = "";
    let skip = 0;
    let get = 5;

    if (payload) {
      search = payload.search;
      skip = payload.skip;
      get = payload.get;
    }

    window.backend
      .getRemoteFiles(search, skip, get)
      .then(({ Result, Count }) => {
        commit("setRemoteFiles", Result);
        console.log(Result, Count);
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
