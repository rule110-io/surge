<template>
  <div class="sidebar">
    <div class="sidebar__controls">
      <div class="sidebar__control" @click="openShowTopicModal">
        <Icon class="sidebar__control-icon" icon="PlusIcon"></Icon>Add New Topic
      </div>
    </div>
    <!-- <div style="display: flex;">
      <input
        style="width: 80px;"
        type="text"
        v-model.trim="topicName"
        placeholder="topic name"
      />
      <button @click="subscribeToTopic(topicName)">subscribe</button>
    </div> -->

    <div class="sidebar__title">
      Subsribed
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
    <Modal :show.sync="showTopicModal">
      <template slot="title">
        Add New Topic
      </template>
      <template slot="body">
        <ControlWrapper title="Topic name">
          <Input
            v-model="topicName"
            theme="light"
            size="md"
            placeholder="Enter topic name here"
          />
        </ControlWrapper>
      </template>
      <template slot="footer">
        <Button theme="text" size="md" @click="closeShowTopicModal"
          >Close</Button
        >
        <Button
          theme="default"
          size="md"
          @click="subscribeAndActivateTopic(topicName)"
          >Add New</Button
        >
      </template>
    </Modal>
  </div>
</template>

<style lang="scss">
@import "./Sidebar.scss";
</style>

<script>
import Icon from "@/components/Icon/Icon";
import Modal from "@/components/Modals/Modal/Modal";
import ControlWrapper from "@/components/Controls/ControlWrapper/ControlWrapper";
import Input from "@/components/Controls/Input/Input";
import Button from "@/components/Button/Button";

import { mapState, mapActions, mapMutations } from "vuex";

export default {
  components: { Icon, Modal, ControlWrapper, Input, Button },
  data: () => {
    return {
      topicName: "",
      showTopicModal: false,
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
    subscribeAndActivateTopic(topic) {
      this.subscribeToTopic(topic);
      this.setRemoteFilesTopic(topic);
      this.closeShowTopicModal();
      this.clearTopicModal();
    },
    openShowTopicModal() {
      this.showTopicModal = true;
    },
    closeShowTopicModal() {
      this.showTopicModal = false;
    },
    clearTopicModal() {
      this.topicName = "";
    },
  },
};
</script>
