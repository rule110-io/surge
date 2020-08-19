const getDefaultState = () => {
  return {
    total: 0,
    online: 0,
  };
};

const state = getDefaultState();

const mutations = {
  addClientStatus(state, event) {
    const { total, online } = event;
    state.total = total;
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
