const getDefaultState = () => {
  return {
    clientStatus: false,
  };
};

const state = getDefaultState();

const mutations = {
  addClientStatus(state, event) {
    state.clientStatus = event;
  },
};

export default {
  namespaced: true,
  state,
  getters: {},
  actions: {},
  mutations,
};
