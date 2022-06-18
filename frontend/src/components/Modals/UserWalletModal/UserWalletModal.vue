<template>
  <Modal
    class="wallet"
    :show.sync="showModal"
    @closeAndClear="closeAndClearModal"
  >
    <template slot="title"> Wallet </template>
    <template slot="body">
      <ModalGrid>
        <ControlWrapper title="Wallet address">
          <div class="text_select text_descr" v-clipboard:copy="walletAddress">
            {{ walletAddress }}
          </div>
        </ControlWrapper>
        <ControlWrapper title="Wallet balance">
          <div class="text_descr">
            <b>{{ walletBalance.toFixed(8) }}</b> NKN
          </div>
        </ControlWrapper>
        <ControlWrapper title="Transaction fee">
          <template slot="descr">
            <div>
              <span> <b>Low:</b> ~{{ lowFee }} NKN</span><br />
              <span> <b>Average:</b> ~{{ avgFee }} NKN</span><br />
              <span> <b>High:</b> ~{{ highFee }} NKN</span>
            </div>
          </template>
          <div class="settings__slider">
            <VueSlider
              class="settings__slider-control"
              v-if="showModal"
              v-model="txFee"
              v-bind="sliderOptions"
            ></VueSlider>
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
@import "./UserWalletModal.scss";
</style>

<script>
import FormMixin from "@/mixins/FormMixin.js";

import Modal from "@/components/Modals/Modal/Modal";
import ControlWrapper from "@/components/Controls/ControlWrapper/ControlWrapper";
import Button from "@/components/Button/Button";
import ModalGrid from "@/components/Modals/ModalGrid/ModalGrid";
import VueSlider from "vue-slider-component";

import axios from "axios";

import {} from "vuex";

export default {
  mixins: [FormMixin],
  components: {
    Modal,
    ControlWrapper,
    Button,
    ModalGrid,
    VueSlider,
  },
  data: () => {
    return {
      walletAddress: "",
      walletBalance: 0,
      txFee: 0,
      avgFee: "loading...",
      lowFee: "loading...",
      maxFee: "loading...",
      downloadPath: "",
      sliderOptions: {
        dotSize: 20,
        height: 2,
        tooltip: false,
        interval: 1,
        marks: {
          0: "Low",
          50: "Average",
          100: "High",
        },
        adsorb: true,
        included: true,
      },
    };
  },
  watch: {
    showModal() {
      this.getAvgTxFee();
      this.getNumWorkers();
      this.getWalletAddress();
      this.getWalletBalance();
      this.getTxFee();
    },
    txFee(newVal) {
      this.setTxFee(newVal);
    },
  },
  mounted() {},
  methods: {
    getAvgTxFee() {
      axios
        .get("https://openapi.nkn.org/api/v1/statistics/avgtxfee")
        .then((resp) => {
          const avgFee = resp.data;

          const feePercent = 0.2;

          const lowFee = avgFee - avgFee * feePercent;
          const highFee = avgFee + avgFee * feePercent;

          this.avgFee = avgFee;
          this.lowFee = lowFee;
          this.highFee = highFee;
        })
        .catch((err) => {
          console.log(err);
        });
    },
    getTxFee() {
      window.go.surge.MiddlewareFunctions.GetTxFee().then((resp) => {
        this.txFee = parseFloat(resp);
      });
    },
    setTxFee(fee) {
      window.go.surge.MiddlewareFunctions.SetTxFee(fee.toString());
    },
    getWalletAddress() {
      window.go.surge.MiddlewareFunctions.GetWalletAddress().then((resp) => {
        this.walletAddress = resp;
      });
    },
    getWalletBalance() {
      window.go.surge.MiddlewareFunctions.GetWalletBalance().then((resp) => {
        this.walletBalance = parseFloat(resp);
      });
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
