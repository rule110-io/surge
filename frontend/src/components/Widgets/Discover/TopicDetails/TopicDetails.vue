<template>
  <div
    class="topic-details"
    v-if="remoteFilesConfig.topicName.length && topicDetails"
  >
    <div class="topic-details__top">
      <div class="topic-details__title selectable">
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
        <div
          v-if="topicDetails.Name !== officialTopicName"
          class="topic-details__stats-more"
          v-on-clickaway="closeDropdown"
        >
          <div class="topic-details__stats-more-btn" @click="toggleDropdown">
            <Icon class="topic-details__stats-more-icon" icon="MoreIcon" />
          </div>

          <Dropdown class="topic-details__dropdown" :open.sync="dropdownOpen">
            <ul class="dropdown__list">
              <li class="dropdown__list-item" @click="unsubscribe">
                Unsubscribe
              </li>
            </ul>
          </Dropdown>
        </div>
      </div>
    </div>
    <div class="topic-details__bot">
      <Input
        class="topic-details__search"
        v-model="searchQuery"
        icon="SearchIcon"
        placeholder="Filter files..."
        @update="remoteSearch"
      />
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
import Input from "@/components/Controls/Input/Input";

import { mixin as clickaway } from "vue-clickaway";

export default {
  mixins: [clickaway],
  components: { ShieldIcon, Icon, Dropdown, Input },
  data: () => {
    return {
      dropdownOpen: false,
      searchQuery: "",
      remoteSearch: () => {},
    };
  },
  computed: {
    ...mapState("topics", ["topicDetails", "officialTopicName"]),
    ...mapState("files", ["remoteFilesConfig"]),
  },
  watch: {
    "remoteFilesConfig.topicName"(newVal) {
      this.getTopicDetails(newVal);
      this.searchQuery = "";
      this.remoteSearch();
    },
  },
  created() {
    this.initRemoteSearch();
  },
  mounted() {},
  methods: {
    ...mapActions({
      getTopicDetails: "topics/getTopicDetails",
      fetchRemoteFiles: "topics/fetchRemoteFiles",
    }),
    ...mapMutations({
      setRemoteFilesConfig: "files/setRemoteFilesConfig",
    }),
    unsubscribe() {
      this.$bus.$emit("openUnsubscribeTopicModal");
      this.closeDropdown();
    },
    openDropdown() {
      this.dropdownOpen = true;
    },
    closeDropdown() {
      this.dropdownOpen = false;
    },
    toggleDropdown() {
      if (this.dropdownOpen) {
        this.closeDropdown();
      } else {
        this.openDropdown();
      }
    },
    initRemoteSearch() {
      this.remoteSearch = this._.debounce(() => {
        let newConfig = Object.assign({}, this.remoteFilesConfig);
        newConfig.skip = 0;
        newConfig.search = this.searchQuery;
        this.setRemoteFilesConfig(newConfig);
        this.fetchRemoteFiles();
      }, 500);
    },
  },
};
</script>
