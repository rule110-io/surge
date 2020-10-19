<template>
  <div class="file-status">
    <template v-if="!file.IsDownloading && !file.IsUploading && !file.IsPaused">
      Ready
    </template>
    <template v-else-if="file.IsUploading">
      Seeding
    </template>
    <template v-else-if="file.IsPaused">
      Paused
    </template>
    <template v-else>
      Downloading
    </template>
  </div>
</template>

<style lang="scss">
@import "./FileStatus.scss";
</style>

<script>
export default {
  props: {
    file: {
      type: Object,
      default: () => {},
    },
  },
  data: () => {
    return {};
  },
  computed: {},
  watch: {},
  mounted() {},
  methods: {},
};
</script>
