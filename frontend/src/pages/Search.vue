<template>
  <div>
    <h1 class="page__title">Explore</h1>
    <div class="table">
      <div class="table__row">
        <div class="table__head">Name & size</div>
        <div class="table__head">Seeder</div>
      </div>
      <div class="table__row" v-for="file in files" :key="file.Filename">
        <div class="table__cell"><FileInfo :file="file" /></div>
        <div class="table__cell">{{ file.Seeder }}</div>
        <div class="table__cell">
          <feather
            class="table__action"
            type="download"
            @click.native="download(file)"
          ></feather>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import FileInfo from "@/components/File/FileInfo/FileInfo";

export default {
  components: {
    FileInfo,
  },
  data: () => {
    return {
      files: [],
    };
  },
  mounted() {
    this.getRemote();
  },
  methods: {
    getRemote() {
      window.backend.getRemoteFiles().then((result) => {
        this.files = result;
        console.log(result);
      });
    },
    download(file) {
      const { Filename, FileSize, Seeder } = file;
      window.backend.downloadFile(Seeder, FileSize, Filename).then((result) => {
        console.log(result);
      });
    },
  },
};
</script>
