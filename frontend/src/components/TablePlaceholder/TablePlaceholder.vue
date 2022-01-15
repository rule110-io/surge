<template>
  <div class="table-placeholder">
    <TablePlaceholderIcon />
    <div class="table-placeholder__title">
      Start by subscribing to channels and browsing files in the Discover
      section
    </div>
    <Button
      @click="openDiscover"
      theme="transparent"
      size="md"
      icon="DiscoverIcon"
      >Discover</Button
    >
  </div>
</template>

<style lang="scss">
@import "./TablePlaceholder.scss";
</style>

<script>
import TablePlaceholderIcon from "@/assets/icons/TablePlaceholderIcon.svg";
import Button from "@/components/Button/Button";

export default {
  components: { TablePlaceholderIcon, Button },
  computed: {},
  data: () => {
    return {};
  },
  methods: {
    openDiscover() {
      this.$router.push("/discover");
    },
  },
};
</script>
