<template>
  <div class="settings">
    <div class="settings__header">
      <div class="settings__header-version">Manage your surge - v{{ currentVersion }}</div>
      <div class="settings__header-support">
        Community support:
        <a
          @click="openLink('https://discord.gg/Zg3U5cb')"
          class="settings__header-link"
          href="#"
          target="_blank"
        >https://discord.gg/Zg3U5cb</a>
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
        @click="openLink('https://github.com/rule110-io/surge-ui/releases')"
        class="settings__item-link"
        href="#"
        target="_"
      >Get latest</a>
    </div>

    <div class="settings__item">
      <div class="settings__item-left">
        <feather class="settings__item-icon" type="link"></feather>
        <span class="settings__item-title">Official website</span>
      </div>
      <a class="settings__item-link" href="#" target="_">surge.io</a>
    </div>

    <div class="settings__item">
      <div class="settings__item-left">
        <feather class="settings__item-icon" type="book"></feather>
        <span class="settings__item-title">Guide tour</span>
      </div>
      <div class="settings__item-link" @click="startTour">start</div>
    </div>

    <div class="settings__item">
      <div class="settings__item-left">
        <feather class="settings__item-icon" type="file-text"></feather>
        <span class="settings__item-title">Surge logs</span>
      </div>
      <div class="settings__item-link" @click="openLog">open</div>
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
    startTour() {
      this.$router.push("search");
      this.$store.commit("tour/setTour", "true");
      this.$tours["myTour"].start();
    },
    openLog() {
      window.backend.openLog().then(() => {});
    },
    openLink(Link) {
      window.backend.openLink(Link).then(() => {});
    },
  },
};
</script>
