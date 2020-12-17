<template>
  <div class="page">
    <div class="page__results" id="files_table">
      <h1 class="page__title">File Transfers</h1>
      <div class="table">
        <div class="table__row">
          <div class="table__head" style="width: calc(100% - 646px);">
            File
          </div>
          <div
            class="table__head"
            style="width: 110px; justify-content: center;"
          >
            Down
          </div>
          <div
            class="table__head"
            style="width: 110px; justify-content: center;"
          >
            Up
          </div>
          <div
            class="table__head"
            style="width: 156px; justify-content: center;"
          >
            Status
          </div>
          <div
            class="table__head"
            style="width: 120px; justify-content: center;"
          >
            Remaining
          </div>
          <div class="table__head" style="width: 70px;">Seeds</div>
          <div class="table__head" style="width: 80px;"></div>
        </div>
        <TablePlaceholder v-if="!localFiles.length" type="transfer" />
        <template v-else>
          <div
            class="table__row"
            v-for="file in localFiles"
            :key="file.FileName"
          >
            <div class="table__cell" style="width: calc(100% - 646px);">
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
            <div class="table__cell" style="width: 156px;">
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
    };
  },
  computed: {
    ...mapState("files", ["localFiles", "localCount"]),
  },
  mounted() {
    this.$store.dispatch("files/fetchLocalFiles");
  },
  methods: {
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
