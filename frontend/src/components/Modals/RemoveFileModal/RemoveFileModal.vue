le i<template>
  <div class="modal__wrapper" :class="isOpen ? 'modal__wrapper_open' : null">
    <div class="modal">
      <h2 class="modal__title">
        Delete File
      </h2>
      <p class="modal__descr">
        Attention! You are going to remove <b>{{ file.FileName }}</b> from Surge
        - are you sure?
      </p>
      <div class="modal__footer modal__footer_end">
        <div>
          <Checkbox
            class="modal__checkbox_footer"
            name="fromDisk"
            :value="fromDisk"
            @change="changeFromDisk"
          >
            Delete from disk
          </Checkbox>
        </div>

        <div class="modal__footer-controls">
          <Button :click="closeModal" theme="default"> Cancel </Button>
          <Button :click="removeFile" theme="error">
            Delete
          </Button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Button from "@/components/Button/Button";
import Checkbox from "@/components/Controls/Checkbox/Checkbox.vue";

export default {
  components: { Button, Checkbox },
  props: {
    open: {
      type: Boolean,
      default: false,
    },
    file: {
      type: Object,
      default: () => {},
    },
  },
  data() {
    return {
      fromDisk: false,
    };
  },
  computed: {
    isOpen() {
      return this.open;
    },
  },
  methods: {
    changeFromDisk() {
      this.fromDisk = !this.fromDisk;
    },
    closeModal() {
      this.$emit("toggleRemoveFileModal", false);
    },
    removeFile() {
      window.backend.removeFile(this.file.FileHash, this.fromDisk).then(() => {
        this.$store.dispatch("files/fetchLocalFiles");
        this.closeModal();
      });
    },
  },
};
</script>
