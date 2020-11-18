<template>
  <div class="page">
    <div class="page__results" id="search_results">
      <h1 class="page__title">Search: {{ remoteFilesConfig.search }}</h1>
      <div class="table">
        <div class="table__row">
          <div class="table__head">Name & size</div>
          <div class="table__head">Chunks</div>
          <div class="table__head">File Hash</div>
          <div class="table__head">Seeds</div>
          <div class="table__head">Source</div>
        </div>
        <TablePlaceholder v-if="!remoteFiles" type="search" />
        <template v-else>
          <div
            class="table__row"
            v-for="file in remoteFiles"
            :key="file.FileName"
          >
            <div class="table__cell">
              <FileInfo :file="file" :max="true" :icon="false" />
            </div>
            <div class="table__cell text_align_center">
              {{ file.NumChunks }}
            </div>
            <div class="table__cell">
              <FileHash :hash="file.FileHash" />
            </div>
            <div class="table__cell text_align_center">
              {{ file.SeederCount }}
            </div>
            <div class="table__cell">
              <FileSeeders :seeders="file.Seeders" />
            </div>
            <div class="table__cell text_align_right">
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
        </template>
      </div>
    </div>
  </div>
</template>
<script>
import { mapState } from "vuex";

import FileInfo from "@/components/File/FileInfo/FileInfo";
import FileHash from "@/components/File/FileHash/FileHash";
import Pagination from "@/components/Pagination/Pagination";
import TablePlaceholder from "@/components/TablePlaceholder/TablePlaceholder";

export default {
  name: "search",
  components: {
    FileInfo,
    FileHash,
    Pagination,
    TablePlaceholder,
  },
  data: () => {
    return {};
  },
  computed: {
    ...mapState("files", ["remoteFiles", "localFiles", "remoteFilesConfig"]),
  },
  mounted() {
    this.$store.dispatch("files/fetchRemoteFiles");
  },
  methods: {
    download(hash) {
      window.backend.downloadFile(hash).then(() => {
        this.$store.dispatch("files/fetchLocalFiles");
        this.$store.dispatch("files/fetchRemoteFiles");
      });
    },
  },
};
</script>
