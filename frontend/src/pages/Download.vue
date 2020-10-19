<template>
  <div class="page">
    <div class="page__results" id="files_table">
      <h1 class="page__title">File Transfers</h1>
      <div class="table">
        <div class="table__row">
          <div class="table__head">File</div>
          <div class="table__head">Speed</div>
          <div class="table__head text_align_center">Chunks</div>
          <div class="table__head text_align_center">Remaining</div>
          <div class="table__head text_align_center">Status</div>
        </div>
        <div class="table__row" v-for="file in localFiles" :key="file.FileName">
          <div class="table__cell">
            <FileInfo :file="file" />
          </div>
          <div class="table__cell">
            <FileSpeed :file="file" />
          </div>
          <div class="table__cell">
            <FileChunks :file="file" />
          </div>
          <div class="table__cell">
            <FileTime :file="file" />
          </div>
          <div class="table__cell">
            <FileStatus :file="file" />
          </div>
        </div>
        <Pagination
          dispatcher="files/fetchLocalFiles"
          filesConfig="localFilesConfig"
          filePages="localPages"
          commit="files/setLocalFilesConfig"
        />
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
import FileSpeed from "@/components/File/FileSpeed/FileSpeed";
import FileStatus from "@/components/File/FileStatus/FileStatus";
import Pagination from "@/components/Pagination/Pagination";
import RemoveFileModal from "@/components/Modals/RemoveFileModal/RemoveFileModal";

export default {
  name: "download",
  components: {
    FileInfo,
    FileChunks,
    FileTime,
    FileSpeed,
    FileStatus,
    Pagination,
    RemoveFileModal,
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
  },
};
</script>
