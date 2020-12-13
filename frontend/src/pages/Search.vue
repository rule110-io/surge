<template>
  <div class="page">
    <div class="page__results" id="search_results">
      <h1 class="page__title">Remote Files</h1>
      <div class="table">
        <div class="table__row">
          <div
            v-for="header in headers"
            :key="header.title"
            class="table__head"
            :class="[
              header.sortable ? 'table__head_sortable' : '',
              header.orderName === remoteFilesConfig.orderBy
                ? 'table__head_active'
                : '',
            ]"
            @click="header.sortable ? setSorting(header.orderName) : false"
          >
            {{ header.title }}

            <feather
              v-if="header.orderName === remoteFilesConfig.orderBy"
              class="table__head-action"
              :class="!remoteFilesConfig.isDesc ? 'table__head-action_asc' : ''"
              type="chevron-down"
            ></feather>
          </div>
        </div>
        <template v-if="remoteFiles">
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
            <div class="table__cell text_align_center" style="min-width: 85px;">
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
            :count="remoteCount"
          />
        </template>
        <TablePlaceholder
          v-else-if="remoteFilesConfig.search.length > 0"
          type="search"
        />
        <TablePlaceholder v-else type="remote" />
      </div>
    </div>
  </div>
</template>
<script>
import { mapState } from "vuex";

import FileSeeders from "@/components/File/FileSeeders/FileSeeders";
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
    FileSeeders,
  },
  data: () => {
    return {
      headers: [
        {
          title: "Name & size",
          orderName: "FileSize",
          sortable: true,
        },
        {
          title: "Chunks",
          orderName: "",
          sortable: false,
        },
        {
          title: "File Hash",
          orderName: "",
          sortable: false,
        },
        {
          title: "Seeds",
          orderName: "SeederCount",
          sortable: true,
        },
        {
          title: "Source",
          orderName: "",
          sortable: false,
        },
      ],
    };
  },
  computed: {
    ...mapState("files", [
      "remoteFiles",
      "localFiles",
      "remoteFilesConfig",
      "remoteCount",
    ]),
  },
  mounted() {
    this.$store.dispatch("files/fetchRemoteFiles");
  },
  methods: {
    setSorting(orderBy) {
      console.log(orderBy);
      let newConfig = Object.assign({}, this.remoteFilesConfig);
      const currentOrder = newConfig.orderBy;
      const currentIsDesc = newConfig.isDesc;

      newConfig.isDesc = currentOrder === orderBy ? !currentIsDesc : true;
      newConfig.orderBy = orderBy;

      this.$store.commit("files/setRemoteFilesConfig", newConfig);
      this.$store.dispatch("files/fetchRemoteFiles");
    },
    download(hash) {
      window.backend.downloadFile(hash).then(() => {
        this.$store.dispatch("files/fetchLocalFiles");
        this.$store.dispatch("files/fetchRemoteFiles");
        this.$router.replace("/download");
      });
    },
  },
};
</script>
