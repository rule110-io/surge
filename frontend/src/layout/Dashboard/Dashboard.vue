<template>
  <div class="main__wrapper"><Sidebar /> <Dashboard /></div>
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
  },
  methods: {
    enableNotifications() {
      window.wails.Events.On("notificationEvent", (title, text) => {
        const notification = { title, text };
        this.$store.commit("notifications/addNotification", notification);
      });
    },
    enableDownloadEvents() {
      window.wails.Events.On("downloadStatusEvent", (event) => {
        this.$store.commit("downloadEvents/addDownloadEvent", event);
      });
    },
  },
};
</script>
