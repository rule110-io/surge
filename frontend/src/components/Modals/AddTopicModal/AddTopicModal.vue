<template>
  <Modal :show.sync="showModal" @closeAndClear="closeAndClearModal">
    <template slot="title"> Add New Topic </template>
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
      <Button theme="text" size="md" @click="closeAndClearModal">Close</Button>
      <Button
        theme="default"
        size="md"
        @click="subscribeAndActivateTopic(topicName)"
        :disabled="disabled"
        >Add New</Button
      >
    </template>
  </Modal>
</template>

<script>
import FormMixin from "@/mixins/FormMixin.js";

import Modal from "@/components/Modals/Modal/Modal";
import ControlWrapper from "@/components/Controls/ControlWrapper/ControlWrapper";
import Input from "@/components/Controls/Input/Input";
import Button from "@/components/Button/Button";

import { mapActions, mapMutations } from "vuex";

export default {
  mixins: [FormMixin],
  components: { Modal, ControlWrapper, Input, Button },
  data: () => {
    return {
      topicName: "",
    };
  },
  computed: {
    disabled() {
      return !this.topicName.length;
    },
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
      this.closeModal();
      this.clearModal();
    },
    clearModal() {
      this.topicName = "";
    },
  },
};
</script>
