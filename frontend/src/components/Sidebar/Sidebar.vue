<template>
  <div class="sidebar">
    <div class="sidebar__controls">
      <div class="sidebar__control" @click="openAddTopicModal">
        <Icon class="sidebar__control-icon" icon="PlusIcon"></Icon>Add New
        Channel
      </div>
    </div>
    <div class="sidebar__title">Active</div>

    <div class="sidebar__items">
      <div
        class="sidebar__item"
        v-for="(topic, i) in topicNames"
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
import Icon from "@/components/Icon/Icon";

import { mapState, mapMutations } from "vuex";

export default {
  components: { Icon },
  data: () => {
    return {
      topicName: "",
    };
  },
  computed: {
    ...mapState("topics", ["topics"]),
    ...mapState("files", ["remoteFilesConfig"]),
    topicNames() {
      return this._.map(this.topics, "Name");
    },
  },
  mounted() {
    this.setInitialTopic();
  },
  methods: {
    ...mapMutations({
      setRemoteFilesTopic: "files/setRemoteFilesTopic",
    }),
    setInitialTopic() {
      if (this.topicNames.length > 0) {
        this.setRemoteFilesTopic(this.topicNames[0]);
      }
    },
    openAddTopicModal() {
      this.$bus.$emit("openAddTopicModal");
    },
  },
};
</script>
