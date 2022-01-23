<template>
  <table class="table">
    <thead>
      <tr>
        <th class="text_align_left">name</th>
        <th class="text_align_right">size</th>
        <th class="text_align_left" style="width: 134px">progress</th>
        <th class="text_align_left" style="width: 120px">status</th>
        <th class="text_align_center" style="width: 100px">seeds</th>
        <th class="text_align_right" style="width: 100px">down</th>
        <th class="text_align_right" style="width: 100px">UP</th>
        <th class="text_align_right" style="width: 100px">ETA</th>
        <th></th>
      </tr>
    </thead>
    <tbody>
      <tr
        v-for="file in localFiles"
        :key="file.FileHash"
        class="table__row"
        :class="{
          background_error: file.IsMissing,
          table__row_active: isSelectedFile(file.FileHash),
        }"
        @click.ctrl="updateSelectedFiles(file)"
        @click.meta="updateSelectedFiles(file)"
        @click.exact="addSingleSelectedFile(file)"
      >
        <td class="text_wrap_none" style="max-width: 285px">
          <FileName :file="file" />
        </td>
        <td class="text_align_right">
          <FileSize :file="file" />
        </td>
        <td>
          <FileProgress :file="file" />
        </td>
        <td><FileStatus :file="file" /></td>
        <td class="text_align_center">{{ file.NumSeeders }}</td>
        <td class="text_align_right"><FileDown :file="file" /></td>
        <td class="text_align_right"><FileUp :file="file" /></td>
        <td class="text_align_right"><FileTime :file="file" /></td>
        <td class="text_align_right" style="width: 1px">
          <div style="display: flex; align-items: center">
            <Icon
              icon="FolderIcon"
              @click.native.stop="openFolder(file.FileHash)"
            />
            <FileActions @click.native.stop :file="file" />
          </div>
        </td>
      </tr>
    </tbody>
  </table>
</template>

<style lang="scss">
@import "./TransfersTable";
</style>

<script>
import { mapState, mapActions } from "vuex";

import FileName from "@/components/File/FileName/FileName";
import FileSize from "@/components/File/FileSize/FileSize";
import FileProgress from "@/components/File/FileProgress/FileProgress";
import FileActions from "@/components/File/FileActions/FileActions";
import FileDown from "@/components/File/FileDown/FileDown";
import FileUp from "@/components/File/FileUp/FileUp";
import FileTime from "@/components/File/FileTime/FileTime";
import FileStatus from "@/components/File/FileStatus/FileStatus";

import Icon from "@/components/Icon/Icon";

export default {
  components: {
    FileName,
    FileSize,
    FileProgress,
    Icon,
    FileActions,
    FileUp,
    FileDown,
    FileTime,
    FileStatus,
  },
  data: () => {
    return {};
  },
  computed: {
    ...mapState("files", [
      "localFiles",
      "localCount",
      "localFilesConfig",
      "selectedFiles",
    ]),
  },
  mounted() {
    this.$store.dispatch("files/fetchLocalFiles");
    window.addEventListener("mouseup", this.stopDrag);
  },
  methods: {
    ...mapActions({
      updateSelectedFiles: "files/updateSelectedFiles",
      addSingleSelectedFile: "files/addSingleSelectedFile",
    }),
    setSorting(orderBy) {
      let newConfig = Object.assign({}, this.localFilesConfig);
      const currentOrder = newConfig.orderBy;
      const currentIsDesc = newConfig.isDesc;
      newConfig.isDesc = currentOrder === orderBy ? !currentIsDesc : true;
      newConfig.orderBy = orderBy;
      this.$store.commit("files/setLocalFilesConfig", newConfig);
      this.$store.dispatch("files/fetchLocalFiles");
    },
    openFolder(FileHash) {
      window.go.surge.MiddlewareFunctions.OpenFolder(FileHash).then(() => {});
    },
    isSelectedFile(FileHash) {
      return this._.findIndex(this.selectedFiles, ["FileHash", FileHash]) > -1;
    },
  },
};
</script>
