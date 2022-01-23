<template>
  <div>
    <Sidebar />
    <ContentWrapper type="sidebar">
      <TopicDetails />
      <DiscoverTable />
    </ContentWrapper>
  </div>
</template>
<script>
import ContentWrapper from "@/components/ContentWrapper/ContentWrapper";
import Sidebar from "@/components/Sidebar/Sidebar";
import DiscoverTable from "@/components/Widgets/Discover/DiscoverTable/DiscoverTable";
import TopicDetails from "@/components/Widgets/Discover/TopicDetails/TopicDetails";

export default {
  name: "download",
  components: { ContentWrapper, Sidebar, DiscoverTable, TopicDetails },
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
