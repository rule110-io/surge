<template>
  <ContentWrapper>
    <TransfersHeader />
    <TransfersTable />
  </ContentWrapper>
</template>
<script>
import { mapState } from "vuex";

import ContentWrapper from "@/components/ContentWrapper/ContentWrapper";
import TransfersHeader from "@/components/Widgets/Transfers/TransfersHeader/TransfersHeader";
import TransfersTable from "@/components/Widgets/Transfers/TransfersTable/TransfersTable";

export default {
  name: "download",
  components: { ContentWrapper, TransfersHeader, TransfersTable },
  data: () => {
    return {
      isRemoveFileModal: false,
      activeFile: {},
    };
  },
  computed: {
    ...mapState("files", ["localFiles", "localCount", "localFilesConfig"]),
  },
  mounted() {
    this.$store.dispatch("files/fetchLocalFiles");
  },
  methods: {
    setSorting(orderBy) {
      let newConfig = Object.assign({}, this.localFilesConfig);
      const currentOrder = newConfig.orderBy;
      const currentIsDesc = newConfig.isDesc;

      newConfig.isDesc = currentOrder === orderBy ? !currentIsDesc : true;
      newConfig.orderBy = orderBy;

      this.$store.commit("files/setLocalFilesConfig", newConfig);
      this.$store.dispatch("files/fetchLocalFiles");
    },
    toggleRemoveFileModal(bool) {
      this.isRemoveFileModal = bool;
    },
    removeFile(file) {
      this.activeFile = file;
      this.toggleRemoveFileModal(true);
    },
    openFolder(FileHash) {
      window.backend.MiddlewareFunctions.OpenFolder(FileHash).then(() => {});
    },
  },
};
</script>
