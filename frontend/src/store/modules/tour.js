const getDefaultState = () => {
  return {
    tour: false,
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
    window.go.surge.MiddlewareFunctions.WriteSetting("Tour", "false").then(() => {
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
