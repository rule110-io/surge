<template>
  <div class="network-stats">
    <div class="network-stats__left">
      <div
        v-tooltip="{
          content: 'My NKN Public Key',
          placement: 'top-center',
          offset: 5,
        }"
        class="network-stats__address"
      >
        {{ publicKey }}
      </div>
    </div>
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
    return {
      publicKey: "",
    };
  },
  computed: {
    ...mapState("globalBandwidth", ["statusBundle", "totalDown", "totalUp"]),
  },
  mounted() {
    this.getPublicKey();
  },
  methods: {
    getPublicKey() {
      window.go.surge.MiddlewareFunctions.GetPublicKey().then((res) => {
        this.publicKey = res;
      });
    },
  },
};
</script>
