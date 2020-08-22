const getDefaultState = () => {
  return {
    tour: true,
  };
};

const state = getDefaultState();

const mutations = {
  setTour(state, bool) {
    state.tour = bool == "true";
  },
};

const actions = {
  offTour({ commit }) {
    window.backend.writeSetting("Tour", "false").then(() => {
      commit("setTour", "false");
    });
  },
};

export default {
  namespaced: true,
  state,
  getters: {},
  actions,
  mutations,
};
