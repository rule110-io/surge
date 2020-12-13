<template>
  <div
    class="file-hash text_wrap_none"
    v-tooltip="{
      content: tooltipLabel,
      placement: 'bottom-start',
      offset: 5,
    }"
    v-on:click="copyHash(hash)"
  >
    {{ hash }}
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
      tooltipLabel: "click to copy to clipboard",
    };
  },
  methods: {
    copyHash(hash) {
      const el = document.createElement('textarea');
      el.value = hash;
      document.body.appendChild(el);
      el.select();
      document.execCommand('copy');
      document.body.removeChild(el);
    },
  },
};
</script>
