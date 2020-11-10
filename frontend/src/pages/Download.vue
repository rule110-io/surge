<template>
  <div class="page">
    <div class="page__results" id="files_table">
      <h1 class="page__title">File Transfers</h1>
      <div class="table">
        <div class="table__row">
          <div class="table__head">File</div>
          <div class="table__head text_align_center">Down</div>
          <div class="table__head text_align_center">Up</div>
          <div class="table__head text_align_center">Status</div>
          <div class="table__head">Remaining</div>
          <div class="table__head">Seeds</div>
          <div class="table__head"></div>
        </div>
        <TablePlaceholder v-if="!localFiles" type="transfer" />
        <template v-else>
          <div
            class="table__row"
            v-for="file in localFiles"
            :key="file.FileName"
          >
            <div class="table__cell">
              <FileInfo :file="file" />
            </div>
            <div class="table__cell text_align_center">
              <FileDown :file="file" />
            </div>
            <div class="table__cell text_align_center">
              <FileUp :file="file" />
            </div>
            <div class="table__cell"><FileChunks :file="file" /></div>
            <div class="table__cell">
              <FileTime :file="file" />
            </div>
            <div class="table__cell">
              {{ file.SeederCount }}
            </div>
            <div class="table__cell">
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
    ...mapState("files", ["localFiles"]),
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
