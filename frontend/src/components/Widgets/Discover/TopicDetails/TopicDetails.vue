<template>
  <div class="topic-details" v-if="topicDetails">
    <div class="topic-details__top">
      <div class="topic-details__title">
        <ShieldIcon
          v-if="topicDetails.Name === officialTopicName"
          class="topic-details__title-icon"
        />
        #{{ topicDetails.Name }}
      </div>
      <div class="topic-details__controls"></div>
    </div>
    <div class="topic-details__bot">
      <div class="topic-details__descr">descr</div>
      <div>search</div>
    </div>
  </div>
</template>

<style lang="scss">
@import "./TopicDetails.scss";
</style>

<script>
import { mapState, mapActions } from "vuex";

import ShieldIcon from "@/assets/icons/ShieldIcon.svg";

export default {
  components: { ShieldIcon },
  data: () => {
    return {};
  },
  computed: {
    ...mapState("topics", ["topicDetails", "officialTopicName"]),
    ...mapState("files", ["remoteFilesConfig"]),
  },
  watch: {
    "remoteFilesConfig.topicName"(newVal) {
      this.getTopicDetails(newVal);
    },
  },
  mounted() {},
  methods: {
    ...mapActions({
      getTopicDetails: "topics/getTopicDetails",
    }),
  },
};
</script>
