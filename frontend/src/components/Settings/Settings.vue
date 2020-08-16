<template>
  <div class="settings">
    <div class="settings__header">
      <div class="settings__header-version">
        Manage your surge - v{{ currentVersion }}
      </div>
      <div class="settings__header-support">
        Community support:
        <a
          class="settings__header-link"
          href="https://discord.gg/Zg3U5cb"
          target="_blank"
          >https://discord.gg/Zg3U5cb</a
        >
      </div>
    </div>
    <div class="settings__item">
      <div class="settings__item-left">
        <feather class="settings__item-icon" type="sliders"></feather>
        <span class="settings__item-title">Dark Theme</span>
      </div>
      <Switcher name="theme" :value="darkTheme" @change="changeTheme" />
    </div>

    <div class="settings__item">
      <div class="settings__item-left">
        <feather class="settings__item-icon" type="github"></feather>
        <span class="settings__item-title">releases</span>
      </div>
      <div v-if="!isNewVersion">Latest version installed</div>
      <a
        v-else
        class="settings__item-link"
        href="https://github.com/rule110-io/surge-ui/releases"
        target="_"
        >Get latest</a
      >
    </div>

    <div class="settings__item">
      <div class="settings__item-left">
        <feather class="settings__item-icon" type="link"></feather>
        <span class="settings__item-title">Official website</span>
      </div>
      <a class="settings__item-link" href="#" target="_">surge.io</a>
    </div>
  </div>
</template>

<style lang="scss">
@import "./Settings.scss";
</style>

<script>
import { mapState, mapGetters } from "vuex";

import Switcher from "@/components/Controls/Switcher/Switcher.vue";

export default {
  components: { Switcher },
  data: () => {
    return {};
  },
  computed: {
    ...mapGetters({ darkTheme: "darkTheme/getDarkTheme" }),
    ...mapState("version", ["currentVersion", "remoteVersion", "isNewVersion"]),
  },
  methods: {
    changeTheme() {
      this.$store.dispatch("darkTheme/toggleDarkTheme");
    },
  },
};
</script>
