<template>
  <div :class="`file-avatar file-avatar_type_${seeder}`" ref="avatar"></div>
</template>

<style lang="scss">
@import "./FileAvatar.scss";
</style>

<script>
import { createIcon } from "@download/blockies";

export default {
  components: {},
  props: {
    seeder: {
      type: String,
      default: "",
    },
  },
  data: () => {
    return {};
  },
  computed: {},
  mounted() {
    this.generateIcon();
  },
  methods: {
    generateIcon() {
      const avatar = createIcon({
        // All options are optional
        seed: this.seeder, // seed used to generate icon data, default: random
        bgcolor: "#fff", // choose a different background color, default: white
        color: "#02d2b3",
        size: 32, // width/height of the icon in blocks, default: 10
        scale: 2, // width/height of each block in pixels, default: 5
      });

      this.$refs.avatar.appendChild(avatar);
    },
  },
};
</script>
