const getDefaultState = () => {
  return {
    online: 0,
  };
};

const state = getDefaultState();

const mutations = {
  addClientStatus(state, event) {
    const { online } = event;
    state.online = online === 0 ? 1 : online;
  },
};

export default {
  namespaced: true,
  state,
  getters: {},
  actions: {},
  mutations,
};
