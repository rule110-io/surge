<template>
  <div class="network-stats">
    <div class="network-stats__right">
      <div class="network-stats__item text_wrap_none">
        Download: {{ totalDown | prettyBytes(1) }}/s
      </div>
      <div class="network-stats__item text_wrap_none">
        Upload: {{ totalUp | prettyBytes(1) }}/s
      </div>
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
    ...mapState("globalBandwidth", ["statusBundle", "totalDown", "totalUp"]),
  },
  methods: {},
};
</script>
