<template>
  <table>
    <thead>
      <tr>
        <th class="text_align_left">name</th>
        <th class="text_align_right">size</th>
        <th class="text_align_left">progress</th>
        <th class="text_align_left">status</th>
        <th class="text_align_right">seeds</th>
        <th class="text_align_right">down</th>
        <th class="text_align_right">UP</th>
        <th class="text_align_right">ETA</th>
        <th></th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="file in localFiles" :key="file.FileHash">
        <td>
          <FileName :file="file" />
        </td>
        <td class="text_align_right">
          <FileSize :file="file" />
        </td>
        <td>
          <FileProgress :file="file" />
        </td>
        <td></td>
        <td class="text_align_right"></td>
        <td class="text_align_right"></td>
        <td class="text_align_right"></td>
        <td class="text_align_right"></td>
        <td class="text_align_right" style="width: 1px">
          <div style="display: flex; align-items: center">
            <Icon icon="FolderIcon" @click.native="openFolder(file.FileHash)" />
            <FileActions :file="file" />
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
import { mapState } from "vuex";

import FileName from "@/components/File/FileName/FileName";
import FileSize from "@/components/File/FileSize/FileSize";
import FileProgress from "@/components/File/FileProgress/FileProgress";
import FileActions from "@/components/File/FileActions/FileActions";

import Icon from "@/components/Icon/Icon";

export default {
  components: { FileName, FileSize, FileProgress, Icon, FileActions },
  data: () => {
    return {};
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
    openFolder(FileHash) {
      window.go.surge.MiddlewareFunctions.OpenFolder(FileHash).then(() => {});
    },
  },
};
</script>
