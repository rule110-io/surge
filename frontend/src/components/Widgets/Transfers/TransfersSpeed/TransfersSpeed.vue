<template>
  <TransfersInfoCard class="transfers-spped" :active="fileSpeed" title="Speed">
    <template slot="info"> </template>

    <template slot="body">
      <BandwidthChart :file="lastSelected" />
    </template>
  </TransfersInfoCard>
</template>

<style lang="scss">
@import "./TransfersSpeed";
</style>

<script>
import { mapState } from "vuex";

import TransfersInfoCard from "@/components/Widgets/Transfers/TransfersInfoCard/TransfersInfoCard";
import BandwidthChart from "@/components/BandwidthChart/BandwidthChart";

export default {
  components: { TransfersInfoCard, BandwidthChart },
  computed: {
    ...mapState("files", ["fileSpeed", "selectedFiles"]),
  },
  data: () => {
    return {
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
    },
  },
  mounted() {},
  methods: {},
};
</script>
