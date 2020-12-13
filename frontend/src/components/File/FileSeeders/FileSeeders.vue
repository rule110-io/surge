<template>
  <div class="file-seeders">
    <div
      class="file-seeders__item"
      v-for="seeder in seedersSlice"
      :key="seeder"
      v-tooltip="{
        content: seeder === pubKey ? 'It`s me' : seeder,
        placement: 'bottom-center',
        offset: 0,
      }"
    >
      <FileAvatar :seeder="seeder" type="small" />
    </div>
    <div
      v-if="seedsLeft > 0"
      class="file-seeders__item file-seeders__item_more"
    >
      +{{ seedsLeft }}
    </div>
  </div>
</template>

<style lang="scss">
@import "./FileSeeders.scss";
</style>

<script>
import { mapState } from "vuex";

import FileAvatar from "@/components/File/FileAvatar/FileAvatar";

export default {
  components: { FileAvatar },
  props: {
    seeders: {
      type: Array,
      default: () => [],
    },
  },
  data: () => {
    return {
      count: 4,
    };
  },
  computed: {
    ...mapState("pubKey", ["pubKey"]),
    seedersSlice() {
      return this.seeders ? this.seeders.slice(0, this.count) : [];
    },
    seedsLeft() {
      return this.seeders ? this.seeders.length - this.count : 0;
    },
  },
};
</script>
