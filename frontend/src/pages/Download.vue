<template>
  <div class="page">
    <div class="page__results" id="files_table">
      <h1 class="page__title">File Transfers</h1>
      <div class="table">
        <div class="table__row">
          <div
            v-for="header in headers"
            :key="header.title"
            :style="header.style"
            class="table__head"
            :class="[
              header.sortable ? 'table__head_sortable' : '',
              header.orderName === localFilesConfig.orderBy
                ? 'table__head_active'
                : '',
            ]"
            @click="header.sortable ? setSorting(header.orderName) : false"
          >
            {{ header.title }}

            <feather
              v-if="header.orderName === localFilesConfig.orderBy"
              class="table__head-action"
              :class="!localFilesConfig.isDesc ? 'table__head-action_asc' : ''"
              type="chevron-down"
            ></feather>
          </div>
        </div>
        <template v-if="localFiles.length">
          <div
            class="table__row"
            v-for="file in localFiles"
            :key="file.FileHash"
          >
            <div
              class="table__cell text_wrap_none"
              style="width: calc(100% - 666px);"
            >
              <FileInfo :file="file" />
            </div>
            <div
              class="table__cell"
              style="width: 110px; justify-content: center;"
            >
              <FileDown :file="file" />
            </div>
            <div
              class="table__cell"
              style="width: 110px; justify-content: center;"
            >
              <FileUp :file="file" />
            </div>
            <div
              class="table__cell"
              style="width: 176px; justify-content: center;"
            >
              <FileChunks :file="file" />
            </div>
            <div
              class="table__cell"
              style="width: 120px; justify-content: center;"
            >
              <FileTime :file="file" />
            </div>
            <div class="table__cell" style="width: 70px;">
              {{ file.SeederCount }}
            </div>
            <div
              class="table__cell"
              style="width: 80px; justify-content: flex-end;"
            >
              <feather
                class="table__action table__action_remove"
                type="trash-2"
                @click.native="removeFile(file)"
              ></feather>
              <feather
                class="table__action"
                type="folder"
                @click.native="openFolder(file.FileHash)"
              ></feather>
            </div>
          </div>
          <Pagination
            dispatcher="files/fetchLocalFiles"
            filesConfig="localFilesConfig"
            filePages="localPages"
            commit="files/setLocalFilesConfig"
            :count="localCount"
          />
        </template>
        <TablePlaceholder
          v-else-if="localFilesConfig.search.length > 0"
          type="search"
        />
        <TablePlaceholder v-else type="local" />
      </div>
    </div>
    <RemoveFileModal
      :open="isRemoveFileModal"
      :file="activeFile"
      @toggleRemoveFileModal="toggleRemoveFileModal"
    />
  </div>
</template>
<script>
import { mapState } from "vuex";

import FileInfo from "@/components/File/FileInfo/FileInfo";
import FileChunks from "@/components/File/FileChunks/FileChunks";
import FileTime from "@/components/File/FileTime/FileTime";
import FileUp from "@/components/File/FileUp/FileUp";
import FileDown from "@/components/File/FileDown/FileDown";
import Pagination from "@/components/Pagination/Pagination";
import RemoveFileModal from "@/components/Modals/RemoveFileModal/RemoveFileModal";
import TablePlaceholder from "@/components/TablePlaceholder/TablePlaceholder";

export default {
  name: "download",
  components: {
    FileInfo,
    FileChunks,
    FileTime,
    FileUp,
    FileDown,
    Pagination,
    RemoveFileModal,
    TablePlaceholder,
  },
  data: () => {
    return {
      isRemoveFileModal: false,
      activeFile: {},
      headers: [
        {
          title: "File",
          orderName: "FileName",
          sortable: true,
          style: "width: calc(100% - 666px);",
        },
        {
          title: "Down",
          orderName: "",
          sortable: false,
          style: "width: 110px; justify-content: center;",
        },
        {
          title: "Up",
          orderName: "",
          sortable: false,
          style: "width: 110px; justify-content: center;",
        },
        {
          title: "Status",
          orderName: "",
          sortable: false,
          style: "width: 176px; justify-content: center;",
        },
        {
          title: "Remaining",
          orderName: "",
          sortable: false,
          style: "width: 120px; justify-content: center;",
        },
        {
          title: "Seeds",
          orderName: "",
          sortable: false,
          style: "width: 70px;",
        },
        {
          title: "",
          orderName: "",
          sortable: false,
          style: "width: 80px;",
        },
      ],
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
      window.backend.openFolder(FileHash).then(() => {});
    },
  },
};
</script>
