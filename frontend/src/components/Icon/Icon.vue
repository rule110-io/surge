<template>
  <component v-if="icon" :is="icon" class="icon"></component>
  <span v-else></span>
</template>

<style lang="scss">
@import "./Icon";
</style>

<script>
import TransfersIcon from "@/assets/icons/TransfersIcon.svg";
import DiscoverIcon from "@/assets/icons/DiscoverIcon.svg";
import NotificationsIcon from "@/assets/icons/NotificationsIcon.svg";
import SettingsIcon from "@/assets/icons/SettingsIcon.svg";
import InfoIcon from "@/assets/icons/InfoIcon.svg";
import ChartIcon from "@/assets/icons/ChartIcon.svg";
import CopyIcon from "@/assets/icons/CopyIcon.svg";
import DownloadDefaultIcon from "@/assets/icons/DownloadDefaultIcon.svg";
import PlusIcon from "@/assets/icons/PlusIcon.svg";
import CloseIcon from "@/assets/icons/CloseIcon.svg";
import SelectIcon from "@/assets/icons/SelectIcon.svg";

export default {
  components: {
    TransfersIcon,
    DiscoverIcon,
    NotificationsIcon,
    SettingsIcon,
    InfoIcon,
    ChartIcon,
    CopyIcon,
    DownloadDefaultIcon,
    PlusIcon,
    CloseIcon,
    SelectIcon,
  },
  props: {
    icon: {
      type: String,
      default: "",
    },
  },
};
</script>
