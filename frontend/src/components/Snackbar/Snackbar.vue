<template>
  <notifications
    group="notifications"
    position="bottom center"
    width="500"
    class="snackbar"
    :duration="100000"
    animation-name="fade-bottom"
  >
    <template slot="body" slot-scope="props">
      <div
        class="snackbar__body"
        :class="`snackbar__body_type_${props.item.type}`"
      >
        <div class="snackbar__text">{{ props.item.text }}</div>
        <div class="snackbar__close" @click="props.close">
          <CloseIcon class="snackbar__close-icon" />
        </div>
      </div>
    </template>
  </notifications>
</template>

<style lang="scss">
@import "./Snackbar.scss";
</style>

<script>
import CloseIcon from "@/assets/icons/CloseIcon.svg";

export default {
  components: {
    CloseIcon,
  },
  data() {
    return {};
  },
};
</script>
