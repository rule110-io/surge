export const state = () => ({
  snack: "",
  color: "error",
  timeout: false,
});

const mutations = {
  setSnack(state, payload) {
    state = Object.assign(state, payload);
  },
};

const getters = {
  getSnack(state) {
    return state;
  },
};

const actions = {
  updateSnack({ commit }, snack) {
    commit("setSnack", snack);
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
