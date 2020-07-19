<template>
  <div class="page">
    <h1 class="page__title">Explore</h1>
    <div class="table">
      <div class="table__row">
        <div class="table__head">Name & size</div>
        <div class="table__head">Seeder</div>
      </div>
      <div class="table__row" v-for="file in remoteFiles" :key="file.FileName">
        <div class="table__cell"><FileInfo :file="file" /></div>
        <div class="table__cell">{{ file.Seeder }}</div>
        <div class="table__cell">
          <feather
            v-if="!file.IsTracked"
            class="table__action"
            type="download"
            @click.native="download(file.FileHash)"
          ></feather>
          <feather
            v-if="file.IsTracked"
            class="table__action table__action_active"
            type="check-circle"
          ></feather>
        </div>
      </div>
      <Pagination
        dispatcher="files/fetchRemoteFiles"
        filesConfig="remoteFilesConfig"
        filePages="remotePages"
        commit="files/setRemoteFilesConfig"
      />
    </div>
    <h2 class="page__subtitle">Recent Files</h2>
    <RecentFiles :files="localFiles" />
  </div>
</template>
<script>
import { mapState } from "vuex";

import FileInfo from "@/components/File/FileInfo/FileInfo";
import RecentFiles from "@/components/File/RecentFiles/RecentFiles";
import Pagination from "@/components/Pagination/Pagination";

export default {
  components: {
    FileInfo,
    RecentFiles,
    Pagination,
  },
  data: () => {
    return {};
  },
  computed: {
    ...mapState("files", ["remoteFiles", "localFiles"]),
  },
  mounted() {},
  methods: {
    download(hash) {
      window.backend.downloadFile(hash).then(() => {
        this.$store.dispatch("files/fetchLocalFiles");
      });
    },
  },
};
</script>
