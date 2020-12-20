<template>
  <div class="table-placeholder">
    <div class="table-placeholder__title">
      <template v-if="type === 'transfer'">
        No files are currently being transferred. <br />
        Share something with the `+` icon below, or use search to browse.
      </template>
      <template v-if="type === 'search'">
        No search results for '{{ remoteFilesConfig.search }}'. <br />
        But you can share something with the `+` icon below.
      </template>
      <template v-if="type === 'local'">
        Currently there are no local files. <br />
        Share something with the `+` icon below, or use search to browse.
      </template>
      <template v-if="type === 'remote'">
        Currently there are no remote files. <br />
        Share something with the `+` icon below, or use search to browse.
      </template>
    </div>
  </div>
</template>

<style lang="scss">
@import "./TablePlaceholder.scss";
</style>

<script>
import { mapState } from "vuex";

export default {
  props: {
    type: {
      type: String,
      default: "remote",
    },
  },
  components: {},
  computed: {
    ...mapState("files", ["remoteFilesConfig"]),
  },
  data: () => {
    return {};
  },
  methods: {},
};
</script>
