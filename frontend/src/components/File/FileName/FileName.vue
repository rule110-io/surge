<template>
  <div class="file-name">
    <component :is="icon" class="file-name__icon"></component>
    <div
      class="file-name__title text_wrap_none"
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
@import "./FileName.scss";
</style>

<script>
import DownloadIcon from "@/assets/icons/DownloadIcon.svg";
import UploadIcon from "@/assets/icons/UploadIcon.svg";
import CheckIcon from "@/assets/icons/CheckIcon.svg";

export default {
  components: { DownloadIcon, UploadIcon, CheckIcon },
  props: {
    file: {
      type: Object,
      default: () => {},
    },
  },
  computed: {
    icon() {
      if (this.file.IsDownloading) {
        return "DownloadIcon";
      } else if (this.file.IsUploading) {
        return "UploadIcon";
      } else {
        return "CheckIcon";
      }
    },
  },
  data: () => {
    return {};
  },
  methods: {},
};
</script>
