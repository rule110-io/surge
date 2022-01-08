<template>
  <Modal
    class="settings"
    :show.sync="showModal"
    @closeAndClear="closeAndClearModal"
  >
    <template slot="title"> Settings </template>
    <template slot="body">
      <ModalGrid>
        <ControlWrapper title="Download Folder">
          <div class="settings__path">
            <div class="settings__path-control">
              {{ downloadPath || "No path specified" }}
            </div>
            <PathIcon
              @click.native="setDownloadFolder"
              class="settings__path-icon"
            />
          </div>
        </ControlWrapper>
        <ControlWrapper title="Workers per client">
          <template slot="descr">
            <div>Workers per client descr</div>
          </template>
          <div class="settings__slider">
            <VueSlider
              class="settings__slider-control"
              v-if="showModal"
              v-model="numWorkers"
              :min="minWorkers"
              :max="maxWorkers"
              v-bind="sliderOptions"
            ></VueSlider>
            <Input
              class="settings__slider-input"
              v-model="numWorkers"
              type="number"
              theme="light"
              size="md"
              placeholder="1"
            />
          </div>
        </ControlWrapper>

        <ControlWrapper title="Number of NKN multiclients">
          <template slot="descr">
            <div>Number of NKN multiclients descr</div>
          </template>
          <div class="settings__slider">
            <VueSlider
              class="settings__slider-control"
              v-if="showModal"
              v-model="numClients"
              :min="minClients"
              :max="maxClients"
              v-bind="sliderOptions"
            ></VueSlider>
            <Input
              class="settings__slider-input"
              v-model="numClients"
              type="number"
              theme="light"
              size="md"
              placeholder="1"
            />
          </div>
        </ControlWrapper>
      </ModalGrid>
    </template>
    <template slot="footer">
      <Button theme="default" size="md" @click="closeModal">Close</Button>
    </template>
  </Modal>
</template>

<style lang="scss">
@import "./SettingsModal.scss";
</style>

<script>
import FormMixin from "@/mixins/FormMixin.js";

import Modal from "@/components/Modals/Modal/Modal";
import ControlWrapper from "@/components/Controls/ControlWrapper/ControlWrapper";
import Input from "@/components/Controls/Input/Input";
import Button from "@/components/Button/Button";
import ModalGrid from "@/components/Modals/ModalGrid/ModalGrid";
import VueSlider from "vue-slider-component";
import PathIcon from "@/assets/icons/PathIcon.svg";

import {} from "vuex";

export default {
  mixins: [FormMixin],
  components: {
    Modal,
    ControlWrapper,
    Input,
    Button,
    ModalGrid,
    VueSlider,
    PathIcon,
  },
  data: () => {
    return {
      numWorkers: 8,
      minWorkers: 1,
      maxWorkers: 12,
      numClients: 8,
      minClients: 1,
      maxClients: 8,
      downloadPath: "",
      sliderOptions: {
        dotSize: 20,
        height: 2,
        tooltip: false,
      },
    };
  },
  watch: {
    showModal() {
      this.getDownloadPath();
      this.getNumClients();
      this.getNumWorkers();
    },
    numWorkers(newVal) {
      if (newVal > this.maxWorkers) {
        this.numWorkers = this.maxWorkers;
      }

      if (newVal < this.minWorkers) {
        this.numWorkers = this.minWorkers;
      }

      this.writeSetting("numWorkers", this.numWorkers);
    },
    numClients(newVal) {
      if (newVal > this.maxClients) {
        this.numClients = this.maxClients;
      }

      if (newVal < this.minClients) {
        this.numClients = this.minClients;
      }

      this.writeSetting("numClients", this.numClients);
    },
  },
  mounted() {},
  methods: {
    writeSetting(k, v) {
      window.go.surge.MiddlewareFunctions.WriteSetting(`${k}`, `${v}`);
    },
    setDownloadFolder() {
      window.go.surge.MiddlewareFunctions.SetDownloadFolder();
    },
    getDownloadPath() {
      window.go.surge.MiddlewareFunctions.ReadSetting("downloadFolder").then(
        (res) => {
          this.downloadPath = res;
          console.log(res);
        }
      );
    },
    getNumClients() {
      window.go.surge.MiddlewareFunctions.ReadSetting("numClients").then(
        (res) => {
          this.numClients = Number(res);
        }
      );
    },
    getNumWorkers() {
      window.go.surge.MiddlewareFunctions.ReadSetting("numWorkers").then(
        (res) => {
          this.numWorkers = Number(res);
        }
      );
    },
  },
};
</script>
