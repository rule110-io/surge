<template>
  <div
    :class="['notifications', open ? 'notifications_visible' : '']"
    v-on-clickaway="closeNotifications"
  >
    <div class="notifications__wrapper">
      <div class="notifications__header">
        <div class="notifications__header-title">Notifications</div>
        <div class="notifications__header-clear" @click="clearNotifications">
          Clear all
        </div>
      </div>
      <div class="notifications__body">
        <template v-if="!notifications.length">
          <div class="notifications__item-text">You have no notifications</div>
        </template>
        <template v-else>
          <div
            v-for="(item, i) in notifications"
            :key="i"
            class="notifications__item"
          >
            <div class="notifications__item-header">
              <div class="notifications__item-title">{{ item.title }}</div>
              <div class="notifications__item-time">
                {{ $moment(item.timestamp * 1000).fromNow() }}
              </div>
            </div>

            <div class="notifications__item-text">{{ item.text }}</div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<style lang="scss">
@import "./Notifications.scss";
</style>

<script>
import { mixin as clickaway } from "vue-clickaway";

import { mapState } from "vuex";

export default {
  components: {},
  mixins: [clickaway],
  data: () => {
    return {};
  },
  computed: {
    ...mapState("notifications", ["counter", "open", "notifications"]),
  },
  mounted() {},
  methods: {
    closeNotifications() {
      if (!this.open) return;

      this.$store.commit("notifications/toggleNotifications", false);
    },
    clearNotifications() {
      this.$store.commit("notifications/clearNotifications");
    },
  },
};
</script>
