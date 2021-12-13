<template>
  <TransfersInfoCard
    class="transfers-details"
    :active="fileDetails"
    title="Details"
  >
    <template slot="info">
      <div class="transfers-details__tabs">
        <div
          class="transfers-details__item"
          @click="setActiveTab('General')"
          :class="{ 'transfers-details__item_active': activeTab === 'General' }"
        >
          General
        </div>
        <div
          class="transfers-details__item"
          @click="setActiveTab('Peers')"
          :class="{ 'transfers-details__item_active': activeTab === 'Peers' }"
        >
          Peers
        </div>
      </div>
    </template>

    <template slot="body">
      <div v-show="activeTab === 'General'">
        <FileChunks :file="lastSelected" />
      </div>
    </template>
  </TransfersInfoCard>
</template>

<style lang="scss">
@import "./TransfersDetails";
</style>

<script>
import { mapState } from "vuex";

import TransfersInfoCard from "@/components/Widgets/Transfers/TransfersInfoCard/TransfersInfoCard";
import FileChunks from "@/components/File/FileChunks/FileChunks";

export default {
  components: { TransfersInfoCard, FileChunks },
  computed: {
    ...mapState("files", ["fileDetails", "selectedFiles"]),
  },
  data: () => {
    return {
      activeTab: "General",
      activeFileDetails: null,
      lastSelected: null,
    };
  },

  watch: {
    selectedFiles(newItems) {
      if (!newItems.length) {
        this.lastSelected = null;
        return;
      }

      const lastSelected = newItems[newItems.length - 1];
      this.lastSelected = lastSelected;

      window.go.surge.MiddlewareFunctions.GetFileDetails(
        lastSelected.FileHash
      ).then((resp) => {
        this.activeFileDetails = resp;

        console.log(lastSelected, resp);
      });
    },
  },
  mounted() {},
  methods: {
    setActiveTab(str) {
      this.activeTab = str;
    },
  },
};
</script>
