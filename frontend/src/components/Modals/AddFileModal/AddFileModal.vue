<template>
  <Modal :show.sync="showModal" @closeAndClear="closeAndClearModal">
    <template slot="title"> Add New File </template>
    <template slot="body">
      <ControlWrapper title="Channel name">
        <Select
          v-model="topicName"
          :items="allowedTopics"
          placeholder="Select Channel"
        />
      </ControlWrapper>
    </template>
    <template slot="footer">
      <Button theme="text" size="md" @click="closeAndClearModal">Close</Button>
      <Button
        :disabled="disabled"
        theme="default"
        size="md"
        @click="addFile(topicName)"
        >Upload File</Button
      >
    </template>
  </Modal>
</template>

<script>
import FormMixin from "@/mixins/FormMixin.js";

import Modal from "@/components/Modals/Modal/Modal";
import ControlWrapper from "@/components/Controls/ControlWrapper/ControlWrapper";
import Button from "@/components/Button/Button";
import Select from "@/components/Controls/Select/Select";

import { mapState } from "vuex";

export default {
  mixins: [FormMixin],
  components: { Modal, ControlWrapper, Button, Select },
  data: () => {
    return {
      topicName: "",
    };
  },
  computed: {
    ...mapState("topics", ["topics"]),
    disabled() {
      return !this.topicName.length;
    },
    allowedTopics() {
      return this._.map(
        this._.filter(this.topics, ["Permissions.CanWrite", true]),
        "Name"
      );
    },
  },
  mounted() {},
  methods: {
    addFile(topicName) {
      this.seedFile(topicName);
      this.closeModal();
      this.clearModal();
    },
    seedFile(topicName) {
      window.go.surge.MiddlewareFunctions.SeedFile(topicName).then(() => {
        this.$store.dispatch("files/fetchLocalFiles");
        this.$store.dispatch("files/fetchRemoteFiles");
      });
    },
    clearModal() {
      this.topicName = "";
    },
  },
};
</script>
