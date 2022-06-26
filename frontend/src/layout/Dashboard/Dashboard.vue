<template>
  <div class="main__wrapper">
    <!-- <div class="main__tour" v-if="tour"></div> -->
    <Dashboard />
    <Snackbar />
  </div>
</template>
<script>
import { mapState, mapActions } from "vuex";

import Dashboard from "@/components/Dashboard/Dashboard";
import Snackbar from "@/components/Snackbar/Snackbar";

export default {
  components: {
    Dashboard,
    Snackbar,
  },
  computed: {
    ...mapState("tour", ["tour"]),
    ...mapState("files", ["remoteFilesConfig"]),
  },
  data: () => {
    return {};
  },
  watch: {
    remoteFilesConfig: {
      deep: true,
      handler() {
        this.fetchRemoteFiles();
      },
    },
  },
  destroyed() {
    clearInterval(this.remoteInterval);
    clearInterval(this.localInterval);
  },
  mounted() {
    this.enableNotifications();
    this.enableGlobalBandwidthEvents();
    this.enableErrorEvents();
    this.enableDarkThemeEvent();

    this.fetchTopics();
    this.fetchRemoteFiles();
    this.fetchDarkTheme();
    this.fetchTour();

    this.updateRemoteVersion();

    this.getPublicKey();

    this.getOfficialTopicName();

    this.remoteInterval = setInterval(this.fetchRemoteFiles, 10000);
    this.localInterval = setInterval(this.fetchLocalFiles, 10000);
  },
  methods: {
    ...mapActions({
      getOfficialTopicName: "topics/getOfficialTopicName",
    }),
    getPublicKey() {
      window.go.surge.MiddlewareFunctions.GetPublicKey().then((address) => {
        this.$store.commit("pubKey/setPubKey", address);
      });
    },
    fetchLocalFiles() {
      this.$store.dispatch("files/fetchLocalFiles");
    },
    fetchTopics() {
      this.$store.dispatch("topics/fetchTopics");
    },
    fetchRemoteFiles() {
      this.$store.dispatch("files/fetchRemoteFiles");
    },
    fetchDarkTheme() {
      window.go.surge.MiddlewareFunctions.ReadSetting("DarkMode").then(
        (bool) => {
          this.$store.commit("darkTheme/setDarkTheme", bool);
        }
      );
    },
    fetchTour() {
      window.go.surge.MiddlewareFunctions.ReadSetting("Tour").then((bool) => {
        this.$store.commit("tour/setTour", bool);
      });
    },
    updateRemoteVersion() {
      this.$store.dispatch("version/updateRemoteVersion");
    },
    enableDarkThemeEvent() {
      window.runtime.EventsOn("darkThemeEvent", (bool) => {
        this.$store.commit("darkTheme/setDarkTheme", bool);
      });
    },
    enableNotifications() {
      window.runtime.EventsOn("notificationEvent", (title, text, timestamp) => {
        const notification = { title, text, timestamp };
        this.$store.commit("notifications/addNotification", notification);
      });
    },
    enableErrorEvents() {
      window.runtime.EventsOn("errorEvent", (title, text) => {
        console.log(title);
        this.$notify({
          group: "notifications",
          text: `${title}: ${text}`,
          type: "error",
        });
      });
    },
    enableGlobalBandwidthEvents() {
      window.runtime.EventsOn(
        "globalBandwidthUpdate",
        (statusBundle, totalDown, totalUp) => {
          this.$store.commit("globalBandwidth/addGlobalBandwidth", {
            statusBundle,
            totalDown,
            totalUp,
          });
        }
      );
    },
  },
};
</script>
