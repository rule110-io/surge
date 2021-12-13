<template>
  <div class="transfer-controls">
    <div class="transfer-controls__left">
      <div class="transfer-controls__selected" v-if="selectedFiles.length">
        <span class="transfer-controls__selected_text"
          >{{ selectedFiles.length }} selected</span
        >
        <Icon
          icon="ControlCloseIcon"
          class="transfer-controls__selected-close"
          @click.native="clearSelectedFiles"
        />
      </div>
      <div class="transfer-controls__actions">
        <Icon
          @click.native="setPause(false)"
          icon="ControlPlayIcon"
          class="transfer-controls__actions-icon"
        />
        <Icon
          @click.native="setPause(true)"
          icon="ControlPauseIcon"
          class="transfer-controls__actions-icon"
        />
      </div>
    </div>
    <div class="transfer-controls__right">
      <Button
        @click="toggleFileDetails"
        theme="icon"
        size="md"
        icon="InfoIcon"
        :active="fileDetails"
        >Details</Button
      >
      <Button
        @click="toggleFileSpeed"
        theme="icon"
        size="md"
        icon="ChartIcon"
        :active="fileSpeed"
        >Speed</Button
      >
    </div>
  </div>
</template>

<style lang="scss">
@import "./TransferControls";
</style>

<script>
import { mapState, mapActions } from "vuex";

import Icon from "@/components/Icon/Icon";
import Button from "@/components/Button/Button";

export default {
  components: { Icon, Button },
  data: () => {
    return {};
  },
  computed: {
    ...mapState("files", ["selectedFiles", "fileSpeed", "fileDetails"]),
  },
  watch: {},
  mounted() {},
  methods: {
    ...mapActions({
      clearSelectedFiles: "files/clearSelectedFiles",
      toggleFileSpeed: "files/toggleFileSpeed",
      toggleFileDetails: "files/toggleFileDetails",
    }),
    setPause(bool) {
      const hashes = this._.map(this.selectedFiles, "FileHash");

      window.go.surge.MiddlewareFunctions.SetDownloadPause(hashes, bool).then(
        () => {
          this.$store.dispatch("files/fetchLocalFiles");
        }
      );
    },
  },
};
</script>
