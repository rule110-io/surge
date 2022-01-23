<template>
  <div class="transfers-header">
    <Tabs title="Status" :tabs="filters" :active-tab.sync="activeFilter" />
    <Input
      class="transfers-header__search"
      v-model="searchQuery"
      icon="SearchIcon"
      placeholder="Filter files..."
      @update="localSearch"
    />
  </div>
</template>

<style lang="scss">
@import "./TransfersHeader";
</style>

<script>
import { mapState, mapActions, mapMutations } from "vuex";

import Input from "@/components/Controls/Input/Input";
import Tabs from "@/components/Tabs/Tabs";

export default {
  components: { Input, Tabs },
  data: () => {
    return {
      searchQuery: "",
      filters: ["All", "Downloading", "Seeding", "Completed", "Paused"],
      activeFilter: "",
      localSearch: () => {},
    };
  },
  computed: {
    ...mapState("files", ["localFilesConfig"]),
  },
  watch: {
    activeFilter(newStatus) {
      const filterCode = this._.indexOf(this.filters, newStatus);

      let newConfig = Object.assign({}, this.localFilesConfig);
      newConfig.filter = filterCode;
      this.$store.commit("files/setLocalFilesConfig", newConfig);

      this.$store.dispatch("files/fetchLocalFiles");
    },
  },
  created() {
    this.initLocalSearch();
  },
  mounted() {
    this.activeFilter = "All";
    this.localSearch();
  },
  methods: {
    ...mapActions({
      fetchLocalFiles: "files/fetchLocalFiles",
    }),
    ...mapMutations({
      setLocalFilesConfig: "files/setLocalFilesConfig",
    }),
    initLocalSearch() {
      this.localSearch = this._.debounce(() => {
        let newConfig = Object.assign({}, this.localFilesConfig);
        newConfig.skip = 0;
        newConfig.search = this.searchQuery;
        this.setLocalFilesConfig(newConfig);
        this.fetchLocalFiles();
      }, 500);
    },
  },
};
</script>
