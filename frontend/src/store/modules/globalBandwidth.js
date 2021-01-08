const getDefaultState = () => {
  return {
    statusBundle: [],
    totalDown: 0,
    totalUp: 0,
  };
};

const state = getDefaultState();

const mutations = {
  addGlobalBandwidth(state, event) {
    const { statusBundle, totalDown, totalUp } = event;
    state.statusBundle = statusBundle;
    state.totalDown = totalDown;
    state.totalUp = totalUp;
  },
};

export default {
  namespaced: true,
  state,
  getters: {},
  actions: {},
  mutations,
};
