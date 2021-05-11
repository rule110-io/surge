<template>
  <div>
    <Sidebar />
    <ContentWrapper type="sidebar">
      {{ remoteFiles }}
    </ContentWrapper>
  </div>
</template>
<script>
import { mapState } from "vuex";

import ContentWrapper from "@/components/ContentWrapper/ContentWrapper";
import Sidebar from "@/components/Sidebar/Sidebar";

export default {
  name: "download",
  components: { ContentWrapper, Sidebar },
  data: () => {
    return {
      isRemoveFileModal: false,
      activeFile: {},
    };
  },
  computed: {
    ...mapState("files", ["remoteFiles"]),
  },
  mounted() {
    this.$store.dispatch("files/fetchLocalFiles");
  },
  methods: {},
};
</script>
