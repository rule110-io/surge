const getDefaultState = () => {
  return {
    topics: [],
    topicDetails: null,
    officialTopicName: null,
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
  setTopicDetails(state, topicDetails) {
    state.topicDetails = topicDetails;
  },
  setOfficialTopicName(state, topicName) {
    state.officialTopicName = topicName;
  },
};

const actions = {
  fetchTopics({ commit }) {
    window.go.surge.MiddlewareFunctions.GetTopicSubscriptions().then(
      (topics) => {
        commit("setTopics", topics);
        console.log(topics);
      }
    );
  },
  subscribeToTopic({ dispatch }, topicName) {
    window.go.surge.MiddlewareFunctions.SubscribeToTopic(topicName).then(() => {
      dispatch("fetchTopics");
    });
  },
  unsubscribeFromTopic({ dispatch }, topicName) {
    window.go.surge.MiddlewareFunctions.UnsubscribeFromTopic(topicName).then(
      () => {
        dispatch("fetchTopics");
      }
    );
  },
  getTopicDetails({ commit }, topicName) {
    window.go.surge.MiddlewareFunctions.GetTopicDetails(topicName).then(
      (details) => {
        commit("setTopicDetails", details);
      }
    );
  },
  getOfficialTopicName({ commit }) {
    window.go.surge.MiddlewareFunctions.GetOfficialTopicName().then((name) => {
      commit("setOfficialTopicName", name);
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
