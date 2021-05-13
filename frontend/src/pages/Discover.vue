<template>
  <div>
    <Sidebar />
    <ContentWrapper type="sidebar">
      <DiscoverTable />
    </ContentWrapper>
  </div>
</template>
<script>
import ContentWrapper from "@/components/ContentWrapper/ContentWrapper";
import Sidebar from "@/components/Sidebar/Sidebar";
import DiscoverTable from "@/components/Widgets/Discover/DiscoverTable/DiscoverTable";

export default {
  name: "download",
  components: { ContentWrapper, Sidebar, DiscoverTable },
  data: () => {
    return {
      isRemoveFileModal: false,
      activeFile: {},
    };
  },
  computed: {},
  mounted() {},
  methods: {},
};
</script>
