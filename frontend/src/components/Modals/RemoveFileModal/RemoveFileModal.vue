<template>
  <Modal :show.sync="showModal" @closeAndClear="closeAndClearModal">
    <template slot="title"> Remove File </template>
    <template slot="body">
      <div v-if="activeFile" class="modal__descr modal__item">
        Attention! You are going to remove <b>{{ activeFile.FileName }}</b> from
        Surge - are you sure?
      </div>
      <Checkbox
        class="col_6"
        v-model="fromDisk"
        :value="fromDisk"
        title="Delete from disk"
      />
    </template>
    <template slot="footer">
      <Button theme="text" size="md" @click="closeAndClearModal">Close</Button>
      <Button theme="default" size="md" @click="removeFile">Remove file</Button>
    </template>
  </Modal>
</template>

<script>
import { mapState, mapActions } from "vuex";

import FormMixin from "@/mixins/FormMixin.js";

import Modal from "@/components/Modals/Modal/Modal";
import Button from "@/components/Button/Button";
import Checkbox from "@/components/Controls/Checkbox/Checkbox.vue";

export default {
  mixins: [FormMixin],
  components: { Modal, Button, Checkbox },
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
  },
  mounted() {},
  methods: {
    ...mapActions({
      clearSelectedFiles: "files/clearSelectedFiles",
    }),
    clearModal() {
      this.fromDisk = false;
    },
    removeFile() {
      window.go.surge.MiddlewareFunctions.RemoveFile(
        this.activeFile.FileHash,
        this.fromDisk
      ).then(() => {
        this.clearSelectedFiles();
        let newConfig = Object.assign({}, this.localFilesConfig);
        newConfig.skip = 0;
        this.$store.commit("files/setLocalFilesConfig", newConfig);
        this.$store.dispatch("files/fetchLocalFiles");
        this.closeModal();
        this.clearModal();
      });
    },
  },
};
</script>
