<template>
  <div class="sidebar">
    <div class="sidebar__title">
      Topics
    </div>
    <div style="display: flex;">
      <input
        style="width: 80px;"
        type="text"
        v-model.trim="topicName"
        placeholder="topic name"
      />
      <button @click="subscribeToTopic(topicName)">subscribe</button>
    </div>

    <div class="sidebar__items">
      <div
        class="sidebar__item"
        v-for="(topic, i) in topics"
        :key="i"
        @click="setRemoteFilesTopic(topic)"
        :class="
          remoteFilesConfig.topicName === topic ? 'sidebar__item_active' : null
        "
      >
        {{ topic }}
      </div>
    </div>
  </div>
</template>

<style lang="scss">
@import "./Sidebar.scss";
</style>

<script>
import { mapState, mapActions, mapMutations } from "vuex";

export default {
  components: {},
  data: () => {
    return {
      topicName: "",
    };
  },
  computed: {
    ...mapState("topics", ["topics", "activeTopic"]),
    ...mapState("files", ["remoteFilesConfig"]),
  },
  mounted() {},
  methods: {
    ...mapActions({
      subscribeToTopic: "topics/subscribeToTopic",
    }),
    ...mapMutations({
      setRemoteFilesTopic: "files/setRemoteFilesTopic",
    }),
  },
};
</script>
