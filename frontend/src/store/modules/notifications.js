const getDefaultState = () => {
  return {
    notifications: [],
    open: false,
    counter: 0,
  };
};

const state = getDefaultState();

const mutations = {
  addNotification(state, notification) {
    state.notifications.push(notification);
    state.counter += 1;
  },
  toggleNotification(state, bool) {
    if (bool) {
      state.counter = 0;
    }
    state.open = bool;
  },
};

export default {
  namespaced: true,
  state,
  getters: {},
  actions: {},
  mutations,
};
