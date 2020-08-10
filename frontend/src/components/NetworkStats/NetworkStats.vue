<template>
  <div class="network-stats">
    <div class="network-stats__file" @click="seedFile">
      <feather class="network-stats__file-icon" type="plus"></feather>
    </div>
    <div class="network-stats__item">Output: 300 MB/s</div>
  </div>
</template>

<style lang="scss">
@import "./NetworkStats.scss";
</style>

<script>
export default {
  data: () => {
    return {};
  },
  methods: {
    seedFile() {
      window.backend.seedFile().then(() => {
        this.$store.dispatch("files/fetchLocalFiles");
        this.$store.dispatch("files/fetchRemoteFiles");
      });
    },
  },
};
</script>
