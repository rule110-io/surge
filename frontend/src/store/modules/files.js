const getDefaultState = () => {
  return {
    localFiles: [],
    remoteFiles: [],
    localCount: 0,
    remoteCount: 0,
    localPages: 0,
    remotePages: 0,
    localFilesConfig: {
      filter: 0,
      search: "",
      orderBy: "FileName",
      isDesc: true,
      skip: 0,
      get: 8,
    },
    remoteFilesConfig: {
      topicName: "",
      search: "",
      orderBy: "SeederCount",
      isDesc: true,
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
  setRemoteFilesTopic(state, topicName) {
    state.remoteFilesConfig.topicName = topicName;
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
    const { search, skip, get, orderBy, isDesc, filter } =
      state.localFilesConfig;

    window.go.surge.MiddlewareFunctions.GetLocalFiles(
      search,
      filter,
      orderBy,
      isDesc,
      skip,
      get
    ).then(({ Result, Count }) => {
      commit("setLocalFiles", { Result, Count });
    });
  },
  fetchRemoteFiles({ commit, state }) {
    const { topicName, search, skip, get, orderBy, isDesc } =
      state.remoteFilesConfig;

    window.go.surge.MiddlewareFunctions.GetRemoteFiles(
      topicName,
      search,
      orderBy,
      isDesc,
      skip,
      get
    ).then(({ Result, Count }) => {
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
