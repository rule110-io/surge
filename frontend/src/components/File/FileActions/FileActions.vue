<template>
  <div class="file-actions" v-on-clickaway="closeDropdown">
    <div class="file-actions__btn" @click="toggleDropdown">
      <Icon class="file-actions__icon" icon="MoreIcon" />
    </div>

    <Dropdown class="file-actions__dropdown" :open.sync="dropdownOpen">
      <ul class="dropdown__list">
        <li class="dropdown__list-item" @click="openRemoveFileModal">Remove</li>
      </ul>
    </Dropdown>
  </div>
</template>

<style lang="scss">
@import "./FileActions.scss";
</style>

<script>
import { mixin as clickaway } from "vue-clickaway";
import { mapMutations } from "vuex";

import Icon from "@/components/Icon/Icon";
import Dropdown from "@/components/Dropdown/Dropdown";

export default {
  mixins: [clickaway],
  components: { Icon, Dropdown },
  props: {
    file: {
      type: Object,
      default: () => {},
    },
  },
  data: () => {
    return {
      dropdownOpen: false,
    };
  },
  computed: {},
  mounted() {},
  methods: {
    ...mapMutations({
      setActiveFile: "files/setActiveFile",
    }),
    openDropdown() {
      this.dropdownOpen = true;
    },
    closeDropdown() {
      this.dropdownOpen = false;
    },
    toggleDropdown() {
      if (this.dropdownOpen) {
        this.closeDropdown();
      } else {
        this.openDropdown();
      }
    },
    openRemoveFileModal() {
      this.initActiveFile();
      this.$bus.$emit("openRemoveFileModal");
    },
    initActiveFile() {
      this.setActiveFile(this.file);
      this.closeDropdown();
    },
  },
};
</script>
