<template>
  <div class="main__wrapper">
    <Sidebar />
    <Dashboard />
  </div>
</template>
<script>
import Dashboard from "@/components/Dashboard/Dashboard";
import Sidebar from "@/components/Sidebar/Sidebar";

export default {
  components: {
    Dashboard,
    Sidebar,
  },
  mounted() {
    this.enableNotifications();
    this.enableDownloadEvents();
    this.enableClientStatusUpdate();
    this.enableGlobalBandwidthEvents();
    this.enableErrorEvents();

    this.fetchLocalFiles();
    this.fetchRemoteFiles();

    this.getNumberOfRemoteClient();
  },
  methods: {
    fetchLocalFiles() {
      this.$store.dispatch("files/fetchLocalFiles");
    },
    fetchRemoteFiles() {
      this.$store.dispatch("files/fetchRemoteFiles");
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
        console.log("fileEvent:", event);
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
