<template>
  <div :class="['modal__wrapper', open ? 'modal__wrapper_open' : null]">
    <div class="modal">
      <h2 class="modal__title">
        Start download
      </h2>
      <p class="modal__descr">
        Hey, you are going to start downloading the following surge links:
        {{ links }}
      </p>
      <div class="modal__footer">
        <div class="modal__footer-controls">
          <Button :click="closeModal" theme="default"> Cancel </Button>
          <Button :click="startDownloadMagnetLinks" theme="success">
            Confirm
          </Button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Button from "@/components/Button/Button";

export default {
  components: { Button },
  props: {},
  data: () => {
    return {
      links: [],
      open: false,
    };
  },
  mounted() {
    this.initDownloadEvent();
  },
  methods: {
    closeModal() {
      this.open = false;
    },
    startDownloadMagnetLinks() {
      const links = this.links;

      window.backend.startDownloadMagnetLinks(links).then(() => {
        this.$store.dispatch("files/fetchLocalFiles");
        this.$store.dispatch("files/fetchRemoteFiles");
        this.closeModal();
        this.$router.replace("/download");
      });
    },
    initDownloadEvent() {
      window.wails.Events.On("userEvent", (context, payload) => {
        this.open = true;
        this.links = payload;
      });
    },
  },
};
</script>
