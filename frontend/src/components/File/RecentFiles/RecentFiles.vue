<template>
  <div class="recent">
    <FileCard v-for="file in recent" :key="file.FileName" :file="file" />
  </div>
</template>

<style lang="scss">
@import "./RecentFiles.scss";
</style>

<script>
import FileCard from "@/components/File/FileCard/FileCard";

export default {
  components: { FileCard },
  props: {
    files: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    recent() {
      return this.files.slice(0, 4);
    },
  },
  data() {
    return {};
  },
  methods: {},
};
</script>
