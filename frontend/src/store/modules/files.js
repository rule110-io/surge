const getDefaultState = () => {
  return {
    localFiles: false,
    remoteFiles: false,
    localCount: 0,
    remoteCount: 0,
    localPages: 0,
    remotePages: 0,
  };
};

const state = getDefaultState();

const mutations = {
  setLocalFiles(state, payload) {
    const { Result, Count } = payload;
    state.localFiles = Result;
    state.localCount = Count;
    state.localPages = Math.ceil(Count / 5); // 5 is number of displayed items
  },
  setRemoteFiles(state, payload) {
    const { Result, Count } = payload;
    state.remoteFiles = Result;
    state.remoteCount = Count;
    state.remotePages = Math.ceil(Count / 5); // 5 is number of displayed items
  },
};

const actions = {
  fetchLocalFiles({ commit }) {
    window.backend.getLocalFiles(0, 5).then(({ Result, Count }) => {
      commit("setLocalFiles", { Result, Count });
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
