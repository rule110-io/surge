export default {
  props: {
    openModalEvent: {
      type: String,
      default: "openModal",
    },
  },
  data: function () {
    return {
      showModal: false,
    };
  },
  mounted() {
    this.$bus.$on(this.openModalEvent, this.openModal);
  },
  beforeDestroy() {
    this.$bus.$off(this.openModalEvent);
  },
  methods: {
    clearModal() {},
    openModal() {
      this.showModal = true;
    },
    closeModal() {
      this.showModal = false;
    },
    closeAndClearModal() {
      this.closeModal();
      this.clearModal();
    },
  },
};
