<template>
  <div class="page">
    <div class="page__results" id="search_results">
      <h1 class="page__title page__title_explore">Explore</h1>
      <div class="tab">
        <div
          class="tab__item"
          :class="!isRemote ? 'tab__item_active' : ''"
          @click="isRemote = false"
        >
          Local Files
        </div>
        <div
          class="tab__item"
          :class="isRemote ? 'tab__item_active' : ''"
          @click="isRemote = true"
        >
          Remote Files
        </div>
        <div
          class="tab__marker"
          :class="isRemote ? 'tab__marker_right' : ''"
        ></div>
      </div>
      <div class="table" v-if="isRemote">
        <div class="table__row">
          <div class="table__head">Name & size</div>
          <div class="table__head text_align_center">Seeds</div>
          <div class="table__head">Source</div>
        </div>
        <TablePlaceholder v-if="!remoteFiles" type="remote" />
        <template v-else>
          <div
            class="table__row"
            v-for="file in remoteFiles"
            :key="file.FileName"
          >
            <div class="table__cell">
              <FileInfo :file="file" :full="true" :icon="false" />
            </div>
            <div class="table__cell text_align_center">
              {{ file.SeederCount }}
            </div>
            <div class="table__cell">{{ file.FileHash }}</div>
            <div class="table__cell">
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

      <div class="table" v-else>
        <div class="table__row">
          <div class="table__head">Name & size</div>
          <div class="table__head text_align_center">Seeds</div>
          <div class="table__head">Source</div>
        </div>
        <TablePlaceholder v-if="!localFiles" type="local" />
        <template v-else>
          <div
            class="table__row"
            v-for="file in localFiles"
            :key="file.FileName"
          >
            <div class="table__cell">
              <FileInfo :file="file" :full="true" :icon="false" />
            </div>
            <div class="table__cell text_align_center">
              {{ file.SeederCount }}
            </div>
            <div class="table__cell">{{ file.FileHash }}</div>
            <div class="table__cell">
              <feather
                class="table__action table__action_active"
                type="check-circle"
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
  </div>
</template>
<script>
import { mapState } from "vuex";

import FileInfo from "@/components/File/FileInfo/FileInfo";
import Pagination from "@/components/Pagination/Pagination";
import TablePlaceholder from "@/components/TablePlaceholder/TablePlaceholder";

export default {
  name: "explore",
  components: {
    FileInfo,
    Pagination,
    TablePlaceholder,
  },
  data: () => {
    return {
      isRemote: false,
    };
  },
  computed: {
    ...mapState("files", ["remoteFiles", "localFiles"]),
  },
  mounted() {
    this.$store.dispatch("files/fetchRemoteFiles");
    this.$store.dispatch("files/fetchLocalFiles");
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
