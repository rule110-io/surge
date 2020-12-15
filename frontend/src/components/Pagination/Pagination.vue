<template>
  <div class="pagination">
    <div class="pagination__descr">
      Showing {{ showFrom }} to {{ showTo }} of {{ count }}
    </div>
    <div class="pagination__controls">
      <feather
        :class="[
          'pagination__button',
          isPrev ? 'pagination__button_active' : 'pagination__button_inactive',
        ]"
        type="arrow-left-circle"
        @click.native="decreasePage"
      />
      <feather
        :class="[
          'pagination__button',
          isNext ? 'pagination__button_active' : 'pagination__button_inactive',
        ]"
        type="arrow-right-circle"
        @click.native="increasePage"
      />
    </div>
  </div>
</template>

<style lang="scss">
@import "./Pagination.scss";
</style>

<script>
import { mapState } from "vuex";

export default {
  components: {},
  props: {
    dispatcher: {
      type: String,
      default: "",
      required: true,
    },
    filePages: {
      type: String,
      default: "",
      required: true,
    },
    filesConfig: {
      type: String,
      default: "",
      required: true,
    },
    commit: {
      type: String,
      default: "",
      required: true,
    },
    count: {
      type: Number,
      default: 0,
      required: true,
    },
  },
  data: () => {
    return {};
  },
  computed: {
    ...mapState("files", [
      "remotePages",
      "remoteFilesConfig",
      "localFilesConfig",
      "localPages",
    ]),
    isPrev() {
      return this.config.skip > 0;
    },
    isNext() {
      return (
        this.config.skip <
        this[this.filePages] * this.config.get - this.config.get
      );
    },
    showFrom() {
      return this.config.skip === 0 ? 1 : this.config.skip + 1;
    },
    showTo() {
      return this.showFrom + this.config.get > this.count
        ? this.count
        : this.showFrom + this.config.get - 1;
    },
    config() {
      return this[this.filesConfig];
    },
  },
  methods: {
    decreasePage() {
      if (this.isPrev) {
        let newConfig = Object.assign({}, this.config);
        newConfig.skip -= newConfig.get;

        this.$store.commit(this.commit, newConfig);
        this.$store.dispatch(this.dispatcher);
      }
    },
    increasePage() {
      if (this.isNext) {
        let newConfig = Object.assign({}, this.config);
        newConfig.skip += newConfig.get;

        this.$store.commit(this.commit, newConfig);
        this.$store.dispatch(this.dispatcher);
      }
    },
  },
};
</script>
