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
    state.notifications.unshift(notification);
    state.counter += 1;
  },
  toggleNotifications(state, bool) {
    state.open = bool;
  },
  clearNotifications(state) {
    state.counter = 0;
    state.notifications = [];
  },
};

export default {
  namespaced: true,
  state,
  getters: {},
  actions: {},
  mutations,
};
