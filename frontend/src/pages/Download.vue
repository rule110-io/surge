<template>
  <div class="page">
    <div id="files_table">
      <h1 class="page__title">Files</h1>
      <div class="table">
        <div class="table__row">
          <div class="table__head">Name & size</div>
          <div class="table__head">Speed & status</div>
          <div class="table__head" style="width: 150px;">Time</div>
        </div>
        <div class="table__row" v-for="file in localFiles" :key="file.FileName">
          <div class="table__cell">
            <FileInfo :file="file" />
          </div>
          <div class="table__cell">
            <FileStatus :file="file" />
          </div>
          <div class="table__cell">
            <FileTime :file="file" />
          </div>
          <div class="table__cell text_align_right" style="width: 10%;">
            <feather
              v-if="!file.IsPaused && file.IsDownloading"
              class="table__action"
              type="pause-circle"
              @click.native="pause(file.FileHash)"
            ></feather>
            <feather
              v-if="!file.IsDownloading && !file.IsPaused"
              class="table__action table__action_active"
              type="check-circle"
            ></feather>
            <feather
              v-if="file.IsPaused"
              class="table__action"
              type="play-circle"
              @click.native="continueDownload(file.FileHash)"
            ></feather>
            <feather
              class="table__action table__action_remove"
              type="trash-2"
              @click.native="removeFile(file)"
            ></feather>
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
import FileStatus from "@/components/File/FileStatus/FileStatus";
import FileTime from "@/components/File/FileTime/FileTime";
import Pagination from "@/components/Pagination/Pagination";
import RemoveFileModal from "@/components/Modals/RemoveFileModal/RemoveFileModal";

export default {
  name: "download",
  components: {
    FileInfo,
    FileStatus,
    FileTime,
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
    removeFile(file) {
      this.activeFile = file;
      this.toggleRemoveFileModal(true);
    },
    toggleRemoveFileModal(bool) {
      this.isRemoveFileModal = bool;
    },
    pause(hash) {
      window.backend.setDownloadPause(hash, true).then(() => {
        this.$store.dispatch("files/fetchLocalFiles");
      });
    },
    continueDownload(hash) {
      window.backend.setDownloadPause(hash, false).then(() => {
        this.$store.dispatch("files/fetchLocalFiles");
      });
    },
  },
};
</script>
