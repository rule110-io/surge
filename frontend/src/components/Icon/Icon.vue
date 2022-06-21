<template>
  <component
    v-if="icon"
    :is="icon"
    class="icon"
    :class="{ icon_active: active }"
  ></component>
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
import UsersIcon from "@/assets/icons/UsersIcon.svg";
import FileIcon from "@/assets/icons/FileIcon.svg";
import LockIcon from "@/assets/icons/LockIcon.svg";
import MoreIcon from "@/assets/icons/MoreIcon.svg";
import FolderIcon from "@/assets/icons/FolderIcon.svg";
import ControlCloseIcon from "@/assets/icons/ControlCloseIcon.svg";
import ControlPlayIcon from "@/assets/icons/ControlPlayIcon.svg";
import ControlPauseIcon from "@/assets/icons/ControlPauseIcon.svg";
import WalletIcon from "@/assets/icons/WalletIcon.svg";
import TipIcon from "@/assets/icons/TipIcon.svg";

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
    UsersIcon,
    FileIcon,
    LockIcon,
    MoreIcon,
    FolderIcon,
    ControlCloseIcon,
    ControlPlayIcon,
    ControlPauseIcon,
    WalletIcon,
    TipIcon,
  },
  props: {
    icon: {
      type: String,
      default: "",
    },
    active: {
      type: Boolean,
      default: true,
    },
  },
};
</script>
