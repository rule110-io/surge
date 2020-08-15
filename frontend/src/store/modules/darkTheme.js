const getDefaultState = () => {
  return {
    darkTheme: false,
  };
};

const state = getDefaultState();

const mutations = {
  toggleDarkTheme(state) {
    state.darkTheme = !state.darkTheme;
  },
};

export default {
  namespaced: true,
  state,
  getters: {},
  actions: {},
  mutations,
};
