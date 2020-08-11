<template>
  <div class="network-stats">
    <div class="network-stats__item">
      Clients:
      <span class="network-stats__status" v-if="online === 0">
        <feather class="network-stats__loader" type="loader"></feather>
        Loading...
      </span>
      <template v-else>{{ online }} of {{ total }} connected </template>
    </div>
    <div class="network-stats__file" @click="seedFile">
      <feather class="network-stats__file-icon" type="plus"></feather>
    </div>
    <div class="network-stats__item">
      Avg Speed: {{ totalDown | prettyBytes(1) }}/s |
      {{ totalUp | prettyBytes(1) }}/s
    </div>
  </div>
</template>

<style lang="scss">
@import "./NetworkStats.scss";
</style>

<script>
import { mapState } from "vuex";

export default {
  data: () => {
    return {};
  },
  computed: {
    ...mapState("clientStatus", ["total", "online"]),
    ...mapState("globalBandwidth", ["totalDown", "totalUp"]),
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
