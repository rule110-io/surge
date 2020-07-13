<template>
  <div class="page">
    <h1 class="page__title">Explore</h1>
    <div class="table">
      <div class="table__row">
        <div class="table__head">Name & size</div>
        <div class="table__head">Seeder</div>
      </div>
      <div class="table__row" v-for="file in remoteFiles" :key="file.fileName">
        <div class="table__cell"><FileInfo :file="file" /></div>
        <div class="table__cell">{{ file.seeder }}</div>
        <div class="table__cell">
          <feather
            class="table__action"
            type="download"
            @click.native="download(file)"
          ></feather>
        </div>
      </div>
    </div>
    <h2 class="page__subtitle">Recent Files</h2>
    <RecentFiles :files="localFiles" />
  </div>
</template>
<script>
import FileInfo from "@/components/File/FileInfo/FileInfo";
import RecentFiles from "@/components/File/RecentFiles/RecentFiles";

export default {
  components: {
    FileInfo,
    RecentFiles,
  },
  data: () => {
    return {
      remoteFiles: [],
      localFiles: [],
    };
  },
  mounted() {
    this.getRemote();
    this.getLocal();
  },
  methods: {
    getRemote() {
      window.backend.getRemoteFiles().then((result) => {
        this.remoteFiles = result;
        console.log(result);
      });
    },
    getLocal() {
      window.backend.getLocalFiles().then((result) => {
        this.localFiles = result;
        console.log(result);
      });
    },
    download(file) {
      const { fileName, fileSize, seeder } = file;
      window.backend.downloadFile(seeder, fileSize, fileName).then((result) => {
        console.log(result);
      });
    },
  },
};
</script>
