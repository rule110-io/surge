<template>
  <div class="network-stats">
    <div class="network-stats__item text_wrap_none">
      <span class="network-stats__status" v-if="online === 0">
        <feather class="network-stats__loader" type="loader"></feather>
        Discovering network...
      </span>
      <template v-else>Total clients connected: {{ online }} </template>
    </div>
    <div class="network-stats__file" @click="seedFile">
      <div class="network-stats__file-wrapper">
        <feather class="network-stats__file-icon" type="plus"></feather>
      </div>
    </div>
    <BandwidthChart />

    <div class="network-stats__item text_wrap_none">
      <span class="network-stats__avg">
        Avg Speed: {{ totalDown | prettyBytes(1) }}/s |
        {{ totalUp | prettyBytes(1) }}/s</span
      >
    </div>
  </div>
</template>

<style lang="scss">
@import "./NetworkStats.scss";
</style>

<script>
import BandwidthChart from "@/components/BandwidthChart/BandwidthChart";

import { mapState } from "vuex";

export default {
  components: {
    BandwidthChart,
  },
  data: () => {
    return {};
  },
  computed: {
    ...mapState("clientStatus", ["online"]),
    ...mapState("globalBandwidth", ["statusBundle", "totalDown", "totalUp"]),
  },
  methods: {
    seedFile() {
      window.backend.seedFile().then(() => {
        this.$store.dispatch("files/fetchLocalFiles");
        this.$store.dispatch("files/fetchRemoteFiles");
        this.$router.replace("/download");
      });
    },
  },
};
</script>
