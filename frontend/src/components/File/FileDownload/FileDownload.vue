<template>
  <Icon
    icon="DownloadDefaultIcon"
    class="file-download"
    @click.native="download(hash)"
  />
</template>

<style lang="scss">
@import "./FileDownload.scss";
</style>

<script>
import Icon from "@/components/Icon/Icon";

export default {
  components: { Icon },
  props: {
    hash: {
      type: String,
      default: "",
    },
  },
  data: () => {
    return {
      copied: false,
    };
  },
  methods: {
    download(hash) {
      window.go.surge.MiddlewareFunctions.DownloadFile(hash).then(() => {
        this.$store.dispatch("files/fetchLocalFiles");
        this.$store.dispatch("files/fetchRemoteFiles");
      });
    },
  },
};
</script>
