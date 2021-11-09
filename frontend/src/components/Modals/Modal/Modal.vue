<template>
  <transition name="modal">
    <div class="modal" v-show="show">
      <div class="modal__header">
        <div class="modal__header-title"><slot name="title" /></div>
        <div class="modal__header-close">
          <Icon
            class="modal__header-close-icon"
            @click.native="close"
            icon="CloseIcon"
          ></Icon>
        </div>
      </div>
      <div class="modal__body"><slot name="body" /></div>
      <div class="modal__footer"><slot name="footer" /></div>
    </div>
  </transition>
</template>

<style lang="scss">
@import "./Modal.scss";
</style>

<script>
import Icon from "@/components/Icon/Icon";

export default {
  components: { Icon },
  props: {
    show: {
      type: Boolean,
      required: true,
    },
  },
  data: () => {
    return {};
  },
  computed: {},
  methods: {
    close() {
      this.$emit("update:show", false);
    },
  },
};
</script>
