const getDefaultState = () => {
  return {
    downloadEvent: false,
  };
};

const state = getDefaultState();

const mutations = {
  addDownloadEvent(state, event) {
    state.downloadEvent = event;
  },
};

export default {
  namespaced: true,
  state,
  getters: {},
  actions: {},
  mutations,
};
