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
};
</script>
