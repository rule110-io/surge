const getDefaultState = () => {
  return {
    totalDown: 0,
    totalUp: 0,
  };
};

const state = getDefaultState();

const mutations = {
  addGlobalBandwidth(state, event) {
    const { totalDown, totalUp } = event;
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
