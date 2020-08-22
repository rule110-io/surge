<template>
  <div class="main__wrapper">
    <div class="main__tour" v-if="tour"></div>
    <Sidebar />
    <Dashboard />
    <Tour v-if="tour" />
  </div>
</template>
<script>
import { mapState } from "vuex";

import Dashboard from "@/components/Dashboard/Dashboard";
import Sidebar from "@/components/Sidebar/Sidebar";
import Tour from "@/components/Tour/Tour";

export default {
  components: {
    Dashboard,
    Sidebar,
    Tour,
  },
  computed: {
    ...mapState("tour", ["tour"]),
  },
  data: () => {
    return {};
  },
  mounted() {
    this.enableNotifications();
    this.enableDownloadEvents();
    this.enableClientStatusUpdate();
    this.enableGlobalBandwidthEvents();
    this.enableErrorEvents();
    this.enableDarkThemeEvent();

    this.fetchLocalFiles();
    this.fetchRemoteFiles();
    this.fetchDarkTheme();
    this.fetchTour();

    this.updateRemoteVersion();

    this.getNumberOfRemoteClient();
  },
  methods: {
    fetchLocalFiles() {
      this.$store.dispatch("files/fetchLocalFiles");
    },
    fetchRemoteFiles() {
      this.$store.dispatch("files/fetchRemoteFiles");
    },
    fetchDarkTheme() {
      window.backend.readSetting("DarkMode").then((bool) => {
        this.$store.commit("darkTheme/setDarkTheme", bool);
      });
    },
    fetchTour() {
      window.backend.readSetting("Tour").then((bool) => {
        this.$store.commit("tour/setTour", bool);
      });
    },
    updateRemoteVersion() {
      this.$store.dispatch("version/updateRemoteVersion");
    },
    getNumberOfRemoteClient() {
      window.backend
        .getNumberOfRemoteClient()
        .then(({ NumKnown, NumOnline }) => {
          this.$store.commit("clientStatus/addClientStatus", {
            total: NumKnown,
            online: NumOnline,
          });
        });
    },
    enableDarkThemeEvent() {
      window.wails.Events.On("darkThemeEvent", (bool) => {
        this.$store.commit("darkTheme/setDarkTheme", bool);
      });
    },
    enableNotifications() {
      window.wails.Events.On("notificationEvent", (title, text) => {
        const notification = { title, text };
        this.$store.commit("notifications/addNotification", notification);
      });
    },
    enableErrorEvents() {
      window.wails.Events.On("errorEvent", (title, text) => {
        this.$store.dispatch("snackbar/updateSnack", {
          snack: `${title}: ${text}`,
          color: "error",
          timeout: false,
        });
      });
    },
    enableDownloadEvents() {
      window.wails.Events.On("fileStatusEvent", (event) => {
        this.$store.commit("downloadEvents/addDownloadEvent", event);
      });
    },
    enableGlobalBandwidthEvents() {
      window.wails.Events.On("globalBandwidthUpdate", (totalDown, totalUp) => {
        this.$store.commit("globalBandwidth/addGlobalBandwidth", {
          totalDown,
          totalUp,
        });
      });
    },
    enableClientStatusUpdate() {
      window.wails.Events.On("remoteClientsUpdate", (total, online) => {
        this.$store.commit("clientStatus/addClientStatus", { total, online });
      });
    },
  },
};
</script>
