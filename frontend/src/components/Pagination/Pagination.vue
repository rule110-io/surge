<template>
  <div class="pagination">
    <feather
      class="pagination__button"
      :class="
        isPrev ? 'pagination__button_active' : 'pagination__button_inactive'
      "
      type="arrow-left-circle"
      @click.native="decreasePage"
    />
    <feather
      class="pagination__button"
      :class="
        isNext ? 'pagination__button_active' : 'pagination__button_inactive'
      "
      type="arrow-right-circle"
      @click.native="increasePage"
    />
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
  },
  data() {
    return {
      current: 1,
      skip: 0,
    };
  },
  computed: {
    ...mapState("files", ["remotePages"]),
    isPrev() {
      return this.current > 1;
    },
    isNext() {
      return this.current < this.remotePages;
    },
  },
  watch: {
    current(current) {
      console.log(current);
      const payload = {
        search: "",
        skip: this.skip,
        get: 5,
      };
      this.$store.dispatch("files/fetchRemoteFiles", payload);
    },
  },
  methods: {
    decreasePage() {
      if (this.current > 1) {
        this.current -= 1;
        this.skip -= 5;
      }
    },
    increasePage() {
      if (this.current < this.remotePages) {
        this.current += 1;
        this.skip += 5;
      }
    },
  },
};
</script>
