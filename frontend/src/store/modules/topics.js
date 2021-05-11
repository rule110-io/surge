const getDefaultState = () => {
  return {
    topics: [],
  };
};

const state = getDefaultState();

const getters = {
  getTopics(state) {
    return state.topics;
  },
};

const mutations = {
  setTopics(state, topicsList) {
    state.topics = topicsList;
  },
};

const actions = {
  fetchTopics({ commit }) {
    window.backend.MiddlewareFunctions.GetTopicSubscriptions().then(
      (topics) => {
        commit("setTopics", topics);
      }
    );
  },
  subscribeToTopic({ dispatch }, topicName) {
    window.backend.MiddlewareFunctions.SubscribeToTopic(topicName).then(() => {
      dispatch("fetchTopics");
    });
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
