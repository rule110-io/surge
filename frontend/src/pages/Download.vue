<template>
  <div class="page">
    <h1 class="page__title">Downloading</h1>
    <div class="table">
      <div class="table__row">
        <div class="table__head">Name & size</div>
        <div class="table__head">Speed & status</div>
        <div class="table__head">Time</div>
      </div>
      <div class="table__row" v-for="file in remoteFiles" :key="file.FileName">
        <div class="table__cell"><FileInfo :file="file" /></div>
        <div class="table__cell">
          <FileStatus :file="file" />
        </div>
        <div class="table__cell">
          <FileTime :file="file" />
        </div>
        <div class="table__cell text_align_right" style="width: 10%;">
          <feather
            class="table__action"
            type="pause"
            @click.native="pause(file)"
          ></feather>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import { mapState } from "vuex";

import FileInfo from "@/components/File/FileInfo/FileInfo";
import FileStatus from "@/components/File/FileStatus/FileStatus";
import FileTime from "@/components/File/FileTime/FileTime";

export default {
  components: {
    FileInfo,
    FileStatus,
    FileTime,
  },
  data: () => {
    return {};
  },
  computed: {
    ...mapState("files", ["remoteFiles"]),
  },
  mounted() {},
  methods: {
    pause(file) {
      console.log(file);
    },
  },
};
</script>
