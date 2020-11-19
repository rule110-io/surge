const getDefaultState = () => {
  return {
    pubKey: false,
  };
};

const state = getDefaultState();

const mutations = {
  setPubKey(state, payload) {
    state.pubKey = payload;
  },
};

export default {
  namespaced: true,
  state,
  getters: {},
  actions: {},
  mutations,
};
