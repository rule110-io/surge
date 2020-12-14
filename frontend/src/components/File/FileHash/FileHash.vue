<template>
  <div
    class="file-hash"
    :class="copied ? 'file-hash_active' : null"
    data-label="Copied"
    v-tooltip="{
      content: 'Click to copy to clipboard',
      placement: 'bottom-start',
      offset: 5,
    }"
    v-clipboard:copy="hash"
    v-clipboard:success="onCopy"
  >
    <feather class="file-hash__icon" type="copy"></feather>
    <div class="file-hash__text text_wrap_none">{{ hash }}</div>
  </div>
</template>

<style lang="scss">
@import "./FileHash.scss";
</style>

<script>
export default {
  components: {},
  props: {
    hash: {
      type: String,
      default: "",
    },
  },
  data: () => {
    return {
      copied: false,
    };
  },
  methods: {
    onCopy() {
      this.copied = !this.copied;
    },
  },
};
</script>
