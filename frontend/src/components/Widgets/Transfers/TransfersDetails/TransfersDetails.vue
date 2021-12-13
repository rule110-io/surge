<template>
  <TransfersInfoCard :active="fileDetails" title="Details">
    <template slot="header">
      <div>h</div>
    </template>

    <template slot="body">
      <div>body</div>
    </template>
  </TransfersInfoCard>
</template>

<style lang="scss">
@import "./TransfersDetails";
</style>

<script>
import { mapState } from "vuex";

import TransfersInfoCard from "@/components/Widgets/Transfers/TransfersInfoCard/TransfersInfoCard";

export default {
  components: { TransfersInfoCard },
  computed: {
    ...mapState("files", ["fileDetails"]),
  },
  data: () => {
    return {};
  },

  watch: {},
  mounted() {},
};
</script>
