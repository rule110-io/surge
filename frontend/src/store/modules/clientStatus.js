const getDefaultState = () => {
  return {
    online: 0,
  };
};

const state = getDefaultState();

const mutations = {
  addClientStatus(state, event) {
    state.online = event.online;
  },
};

export default {
  namespaced: true,
  state,
  getters: {},
  actions: {},
  mutations,
};
