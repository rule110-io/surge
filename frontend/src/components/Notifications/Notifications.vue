<template>
  <div class="notifications" :class="open ? 'notifications_visible' : ''">
    <div class="notifications__wrapper">
      <template v-if="!notifications.length">
        <div class="notifications__item">
          <div class="notifications__title">Relax</div>
          <div class="notifications__text">You have no notifications</div>
        </div>
      </template>
      <template v-else>
        <feather
          class="notifications__icon"
          type="check-circle"
          @click.native="closeNotifications"
        ></feather>
        <div
          v-for="(item, i) in notifications"
          :key="i"
          class="notifications__item"
        >
          <div class="notifications__title">{{ item.title }}</div>
          <div class="notifications__text">{{ item.text }}</div>
        </div>
      </template>
    </div>
  </div>
</template>

<style lang="scss">
@import "./Notifications.scss";
</style>

<script>
import { mapState } from "vuex";

export default {
  components: {},
  data: () => {
    return {};
  },
  computed: {
    ...mapState("notifications", ["counter", "open", "notifications"]),
  },
  mounted() {},
  methods: {
    closeNotifications() {
      this.$store.commit("notifications/toggleNotifications", false);
      this.$store.commit("notifications/clearNotifications");
    },
  },
};
</script>
