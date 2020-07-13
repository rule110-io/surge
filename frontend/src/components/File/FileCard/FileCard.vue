<template>
  <div
    class="file-card"
    :style="{ background: `${activeColor}` }"
    @click="openFile(file)"
  >
    <div class="file-card__header">
      <feather class="file-card__icon" type="folder"></feather>
      <div class="file-card__size">{{ file.FileSize | prettyBytes(1) }}</div>
    </div>
    <div class="file-card__footer">
      <div class="file-card__title text_wrap_none">{{ file.Filename }}</div>
      <div class="file-card__progress">Finished — 100%</div>
    </div>
  </div>
</template>

<style lang="scss">
@import "./FileCard.scss";
</style>

<script>
export default {
  props: {
    file: {
      type: Object,
      default: () => {},
    },
  },
  computed: {},
  data() {
    return {
      colors: ["#FEC606", "#2CC990", "#03bf7b", "#8870FF"],
    };
  },
  created() {
    this.activeColor = this.getRandomColor();
  },
  methods: {
    openFile(file) {
      console.log(file);
    },
    getRandomColor() {
      const colors = this.colors;
      return colors[Math.floor(Math.random() * colors.length)];
    },
  },
};
</script>
