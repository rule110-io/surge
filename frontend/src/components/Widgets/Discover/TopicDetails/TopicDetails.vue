<template>
  <div
    class="topic-details"
    v-if="remoteFilesConfig.topicName.length && topicDetails"
  >
    <div class="topic-details__top">
      <div class="topic-details__title">
        <ShieldIcon
          v-if="topicDetails.Name === officialTopicName"
          class="topic-details__title-icon"
        />
        #{{ topicDetails.Name }}
      </div>
      <div class="topic-details__stats">
        <div class="topic-details__stats-item">
          <Icon class="sidebar__stats-icon" :active="false" icon="UsersIcon" />
          <span>{{ topicDetails.Subscribers }}</span>
        </div>
        <div class="topic-details__stats-item">
          <Icon class="sidebar__stats-icon" :active="false" icon="FileIcon" />
          <span>{{ topicDetails.FileCount }}</span>
        </div>
        <div
          class="topic-details__stats-item"
          v-if="!topicDetails.Permissions.CanWrite"
        >
          <Icon class="sidebar__stats-icon" :active="false" icon="LockIcon" />
        </div>
        <div class="topic-details__stats-more" v-on-clickaway="closeDropdown">
          <div class="topic-details__stats-more-btn" @click="openDropdown">
            <Icon class="topic-details__stats-more-icon" icon="MoreIcon" />
          </div>

          <Dropdown class="topic-details__dropdown" :open.sync="dropdownOpen">
            <ul class="dropdown__list">
              <li
                class="dropdown__list-item"
                @click="unsubscribe(topicDetails.Name)"
              >
                Unsubscribe
              </li>
            </ul>
          </Dropdown>
        </div>
      </div>
    </div>
    <div class="topic-details__bot">
      <div>search</div>
    </div>
  </div>
</template>

<style lang="scss">
@import "./TopicDetails.scss";
</style>

<script>
import { mapState, mapActions, mapMutations } from "vuex";

import ShieldIcon from "@/assets/icons/ShieldIcon.svg";
import Icon from "@/components/Icon/Icon";
import Dropdown from "@/components/Dropdown/Dropdown";

import { mixin as clickaway } from "vue-clickaway";

export default {
  mixins: [clickaway],
  components: { ShieldIcon, Icon, Dropdown },
  data: () => {
    return {
      dropdownOpen: false,
    };
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
      unsubscribeFromTopic: "topics/unsubscribeFromTopic",
    }),
    ...mapMutations({
      setRemoteFilesTopic: "files/setRemoteFilesTopic",
    }),
    unsubscribe(topicName) {
      this.unsubscribeFromTopic(topicName);
      this.setRemoteFilesTopic("");
    },
    openDropdown() {
      this.dropdownOpen = true;
    },
    closeDropdown() {
      this.dropdownOpen = false;
    },
  },
};
</script>
