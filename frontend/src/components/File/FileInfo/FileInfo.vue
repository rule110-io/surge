<template>
  <div class="file-info">
    <feather class="file-info__icon" type="file"></feather>
    <div class="file-info__controls">
      <feather
        v-if="!file.IsPaused && file.IsDownloading"
        class="file-info__control"
        type="pause-circle"
        @click.native="pause(file.FileHash)"
      ></feather>
      <feather
        v-if="!file.IsDownloading && !file.IsPaused"
        class="file-info__control"
        type="folder"
        @click.native="openFolder(file.FileHash)"
      ></feather>
      <feather
        v-if="file.IsPaused"
        class="file-info__control"
        type="play-circle"
        @click.native="continueDownload(file.FileHash)"
      ></feather>
      <feather
        class="file-info__control file-info__control_remove"
        type="trash-2"
        @click.native="removeFile(file)"
      ></feather>
    </div>
    <div class="file-info__right">
      <div class="file-info__title">
        <v-clamp autoresize :max-lines="2">{{ file.FileName }}</v-clamp>
      </div>
      <div class="file-info__size">{{ file.FileSize | prettyBytes(1) }}</div>
    </div>
  </div>
</template>

<style lang="scss">
@import "./FileInfo.scss";
</style>

<script>
import VClamp from "vue-clamp";

export default {
  components: {
    VClamp,
  },
  props: {
    file: {
      type: Object,
      default: () => {},
    },
  },
  data: () => {
    return {};
  },
  methods: {
    pause(hash) {
      window.backend.setDownloadPause(hash, true).then(() => {
        this.$store.dispatch("files/fetchLocalFiles");
      });
    },
    removeFile(file) {
      this.activeFile = file;
      this.toggleRemoveFileModal(true);
    },
    continueDownload(hash) {
      window.backend.setDownloadPause(hash, false).then(() => {
        this.$store.dispatch("files/fetchLocalFiles");
      });
    },
    openFolder(FileHash) {
      window.backend.openFolder(FileHash).then(() => {});
    },
  },
};
</script>
