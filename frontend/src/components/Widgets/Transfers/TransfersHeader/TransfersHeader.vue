<template>
  <div class="transfers-header">
    <Tabs title="Status" :tabs="statuses" :active-tab.sync="activeStatus" />
    <Input
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
import Input from "@/components/Controls/Input/Input";
import Tabs from "@/components/Tabs/Tabs";

export default {
  components: { Input, Tabs },
  data: () => {
    return {
      searchQuery: "",
      statuses: ["All", "Downloading", "Seeding", "Completed", "Paused"],
      activeStatus: "All",
    };
  },
};
</script>
