<template>
  <div class="file-info">
    <template v-if="icon">
      <feather
        v-if="file.IsMissing"
        class="file-info__control file-info__control_error"
        type="alert-circle"
        v-tooltip="{
          content: 'Not found on disk!',
          placement: 'bottom-start',
          offset: 5,
        }"
      ></feather>
      <feather
        v-else-if="file.IsHashing"
        class="file-info__control file-info__control_alert"
        type="loader"
        v-tooltip="{
          content: 'Getting ready for sharing, please wait!',
          placement: 'bottom-start',
          offset: 5,
        }"
      ></feather>
      <template v-else>
        <feather
          v-if="!file.IsPaused && file.IsDownloading"
          class="file-info__control"
          type="pause-circle"
          @click.native="pause(file.FileHash)"
        ></feather>
        <feather
          v-if="file.IsPaused"
          class="file-info__control"
          type="play-circle"
          @click.native="continueDownload(file.FileHash)"
        ></feather>
        <feather
          v-if="!file.IsPaused && !file.IsDownloading"
          class="file-info__control file-info__control_active"
          type="check-circle"
        ></feather>
      </template>
    </template>
    <div class="file-info__size text_wrap_none">
      {{ file.FileSize | prettyBytes(1) }}
    </div>
    <div
      class="file-info__title text_wrap_none"
      :class="[
        full ? 'file-info__title_full' : '',
        max ? 'file-info__title_max' : '',
      ]"
      v-tooltip="{
        classes: 'tooltip_left',
        content: file.FileName,
        placement: 'bottom-start',
        offset: 5,
      }"
    >
      {{ file.FileName }}
    </div>
  </div>
</template>

<style lang="scss">
@import "./FileInfo.scss";
</style>

<script>
export default {
  components: {},
  props: {
    file: {
      type: Object,
      default: () => {},
    },
    icon: {
      type: Boolean,
      default: true,
    },
    full: {
      type: Boolean,
      default: false,
    },
    max: {
      type: Boolean,
      default: false,
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
    continueDownload(hash) {
      window.backend.setDownloadPause(hash, false).then(() => {
        this.$store.dispatch("files/fetchLocalFiles");
      });
    },
  },
};
</script>
