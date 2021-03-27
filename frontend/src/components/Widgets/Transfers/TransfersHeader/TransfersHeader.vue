<template>
  <div class="transfers-header">
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

export default {
  components: { Input },
  data: () => {
    return {
      searchQuery: "",
    };
  },
};
</script>
