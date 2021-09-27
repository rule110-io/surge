<template>
  <div class="transfers-header">
    <Tabs title="Status" :tabs="filters" :active-tab.sync="activeFilter" />
    <Input
      class="transfers-header__search"
      :value="searchQuery"
      icon="SearchIcon"
      placeholder="Filter files..."
    />
  </div>
</template>

<style lang="scss">
@import "./TransfersHeader";
</style>

<script>
import { mapState } from "vuex";

import Input from "@/components/Controls/Input/Input";
import Tabs from "@/components/Tabs/Tabs";

export default {
  components: { Input, Tabs },
  data: () => {
    return {
      searchQuery: "",
      filters: ["All", "Downloading", "Seeding", "Completed", "Paused"],
      activeFilter: "",
    };
  },
  computed: {
    ...mapState("files", ["localFilesConfig"]),
  },
  watch: {
    activeFilter(newStatus) {
      const filterCode = this._.indexOf(this.filters, newStatus);

      console.log(filterCode);

      let newConfig = Object.assign({}, this.localFilesConfig);
      newConfig.filter = filterCode;
      this.$store.commit("files/setLocalFilesConfig", newConfig);

      this.$store.dispatch("files/fetchLocalFiles");
    },
  },
  mounted() {
    this.activeFilter = "All";
  },
};
</script>
