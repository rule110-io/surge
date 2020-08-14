<template>
  <div
    :class="[
      'snackbar snackbar_color_' + color,
      show ? 'snackbar_visible' : null,
    ]"
  >
    <p class="snackbar__text">
      {{ message }}
    </p>
    <button class="snackbar__button" @click="close">
      <Feather type="x" class="snackbar__close" />
    </button>
  </div>
</template>

<style lang="scss">
@import "./Snackbar.scss";
</style>

<script>
export default {
  data() {
    return {
      show: false,
      message: "",
      color: "error",
      timeout: false,
    };
  },
  created() {
    this.$store.watch(
      (state) => state.snackbar.snack,
      () => {
        const message = this.$store.state.snackbar.snack;

        if (message !== "") {
          this.show = true;
          this.color = this.$store.state.snackbar.color;
          this.message = message;
          this.timeout = this.$store.state.snackbar.timeout;
          const self = this;

          if (this.timeout === true) {
            setTimeout(function() {
              self.close();
            }, 4000);
          }
        }
      }
    );
  },
  methods: {
    close() {
      this.show = false;
      this.message = "";
      this.color = "error";
      this.timeout = false;

      this.$store.dispatch("snackbar/updateSnack", {
        snack: ``,
        color: "error",
        timeout: false,
      });
    },
  },
};
</script>
