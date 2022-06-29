<template>
  <Modal :show.sync="showModal" @closeAndClear="closeAndClearModal">
    <template slot="title"> Add New Channel </template>
    <template slot="body">
      <ControlWrapper title="Channel name">
        <Input
          v-model="topicName"
          theme="light"
          size="md"
          placeholder="Enter channel name here"
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
      loading: false,
    };
  },
  computed: {
    disabled() {
      if (!this.topicName.length || this.loading) {
        return true;
      } else {
        return false;
      }
    },
  },
  mounted() {},
  methods: {
    ...mapActions({
      fetchTopics: "topics/fetchTopics",
    }),
    ...mapMutations({
      setRemoteFilesTopic: "files/setRemoteFilesTopic",
    }),

    subscribeAndActivateTopic(topic) {
      this.loading = true;
      window.go.surge.MiddlewareFunctions.SubscribeToTopic(topic).finally(
        () => {
          console.log(123);
          this.loading = false;
          this.setRemoteFilesTopic(topic);
          this.fetchTopics();
          this.closeModal();
          this.clearModal();
        }
      );
    },
    clearModal() {
      this.topicName = "";
    },
  },
};
</script>
