<template>
  <div class="network-stats">
    <div class="network-stats__item text_wrap_none">
      <template v-if="online === 0">
        Discovering network...
      </template>
      <template v-else>Total Clients Connected: {{ online }} </template>
    </div>

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
import { mapState } from "vuex";

export default {
  components: {},
  data: () => {
    return {};
  },
  computed: {
    ...mapState("clientStatus", ["online"]),
    ...mapState("globalBandwidth", ["statusBundle", "totalDown", "totalUp"]),
  },
  methods: {},
};
</script>
