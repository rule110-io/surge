<template>
  <table>
    <thead>
      <tr>
        <th class="text_align_left">name</th>
        <th class="text_align_right">size</th>
        <th class="text_align_left">file hash</th>
        <th class="text_align_right">seeds</th>
        <th class="text_align_right"></th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="file in remoteFiles" :key="file.FileHash">
        <td class="text_wrap_none" style="max-width: 300px">
          <FileName name-only :file="file" />
        </td>
        <td class="text_align_right">
          <FileSize :file="file" />
        </td>
        <td class="text_wrap_none" style="max-width: 230px">
          <FileHash :hash="file.FileHash" />
        </td>
        <td class="text_align_right">{{ file.NumSeeders }}</td>
        <td class="text_align_right">
          <div class="discover-table__actions">
            <Icon icon="TipIcon" @click.native="openTipModal(file)" />
            <FileDownload v-if="!file.IsTracked" :hash="file.FileHash" />
            <component
              style="width: 24px"
              v-else
              :is="getFileIcon(file)"
            ></component>
          </div>
        </td>
      </tr>
    </tbody>
  </table>
</template>

<style lang="scss">
@import "./DiscoverTable.scss";
</style>

<script>
import { mapState, mapMutations } from "vuex";

import FileName from "@/components/File/FileName/FileName";
import FileSize from "@/components/File/FileSize/FileSize";
import FileHash from "@/components/File/FileHash/FileHash";
import FileDownload from "@/components/File/FileDownload/FileDownload";

import DownloadIcon from "@/assets/icons/DownloadIcon.svg";
import UploadIcon from "@/assets/icons/UploadIcon.svg";
import CheckIcon from "@/assets/icons/CheckIcon.svg";

import Icon from "@/components/Icon/Icon";

export default {
  components: {
    FileName,
    FileSize,
    FileHash,
    FileDownload,
    DownloadIcon,
    UploadIcon,
    CheckIcon,
    Icon,
  },
  data: () => {
    return {};
  },
  computed: {
    ...mapState("files", ["remoteFiles", "remoteCount", "remoteFilesConfig"]),
  },
  mounted() {
    this.$store.dispatch("files/fetchRemoteFiles");
  },
  methods: {
    ...mapMutations({
      setActiveFile: "files/setActiveFile",
    }),
    openTipModal(file) {
      this.setActiveFile(file);
      this.$bus.$emit("openTipModal");
    },
    setSorting(orderBy) {
      let newConfig = Object.assign({}, this.remoteFilesConfig);
      const currentOrder = newConfig.orderBy;
      const currentIsDesc = newConfig.isDesc;
      newConfig.isDesc = currentOrder === orderBy ? !currentIsDesc : true;
      newConfig.orderBy = orderBy;
      this.$store.commit("files/setRemoteFilesConfig", newConfig);
      this.$store.dispatch("files/fetchRemoteFiles");
    },
    getFileIcon(file) {
      if (file.IsDownloading) {
        return "DownloadIcon";
      } else if (file.IsUploading) {
        return "UploadIcon";
      } else {
        return "CheckIcon";
      }
    },
  },
};
</script>
