<template>
  <Modal :show.sync="showModal" @closeAndClear="closeAndClearModal">
    <template slot="title">Unsubscribe from channel</template>
    <template slot="body">
      <div v-if="topicDetails" class="modal__descr modal__item">
        Attention! You are going to unsubscribe from
        <b>#{{ topicDetails.Name }}</b> channel - are you sure?
      </div>
    </template>
    <template slot="footer">
      <Button theme="text" size="md" @click="closeAndClearModal">Close</Button>
      <Button theme="default" size="md" @click="unsubscribe"
        >Unsubscribe</Button
      >
    </template>
  </Modal>
</template>

<script>
import { mapState, mapMutations, mapActions } from "vuex";

import FormMixin from "@/mixins/FormMixin.js";

import Modal from "@/components/Modals/Modal/Modal";
import Button from "@/components/Button/Button";

export default {
  mixins: [FormMixin],
  components: { Modal, Button },
  props: {
    file: {
      type: Object,
      default: () => {},
    },
  },
  data: () => {
    return {
      fromDisk: false,
    };
  },
  computed: {
    ...mapState("files", ["activeFile", "localFilesConfig"]),
    ...mapState("topics", ["topicDetails"]),
  },
  mounted() {},
  methods: {
    ...mapActions({
      unsubscribeFromTopic: "topics/unsubscribeFromTopic",
    }),
    ...mapMutations({
      setRemoteFilesTopic: "files/setRemoteFilesTopic",
    }),
    unsubscribe() {
      this.unsubscribeFromTopic(this.topicDetails.Name);
      this.setRemoteFilesTopic("");
      this.closeModal();
    },
  },
};
</script>
