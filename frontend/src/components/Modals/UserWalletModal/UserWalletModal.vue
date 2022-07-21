<template>
  <Modal
    class="wallet"
    :show.sync="showModal"
    @closeAndClear="closeAndClearModal"
  >
    <template slot="title"> Wallet </template>
    <template slot="body">
      <Tabs :tabs="tabs" :active-tab.sync="activeTab" class="wallet__tabs" />
      <ModalGrid>
        <template v-if="activeTab === 'General'">
          <ControlWrapper title="Wallet address">
            <div class="text_select text_descr">
              <span v-clipboard:copy="walletAddress">{{ walletAddress }}</span>
              &nbsp;
              <span
                class="text_link"
                @click="openLink(`https://nscan.io/addresses/${walletAddress}`)"
                >(Open in explorer)</span
              >
            </div>
          </ControlWrapper>
          <ControlWrapper title="Wallet balance">
            <div class="text_descr">
              <b>{{ walletBalance.toFixed(8) }}</b> NKN
            </div>
          </ControlWrapper>
          <ControlWrapper title="Channels connecting fee">
            <template slot="descr">
              <div>
                <span> <b>Free:</b> {{ Number(0).toFixed(8) }} NKN</span><br />
                <span> <b>Low:</b> ~{{ lowFee.toFixed(8) }} NKN</span><br />
                <span> <b>Average:</b> ~{{ avgFee.toFixed(8) }} NKN</span><br />
                <span> <b>High:</b> ~{{ highFee.toFixed(8) }} NKN</span>
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
        </template>
        <template v-else>
          <ControlWrapper title="Recipient NKN wallet address">
            <Input
              v-model="transferAddress"
              theme="light"
              size="md"
              placeholder="Enter NKN wallet address"
            />
          </ControlWrapper>
          <ControlWrapper :title="`Transfer amount`">
            <template slot="descr">
              <div>
                Balance: {{ walletBalance.toFixed(8) }}
                <span class="wallet__max" @click="setMaxTransfer">Max</span>
              </div>
            </template>
            <Input
              v-model="transferAmount"
              type="number"
              theme="light"
              size="md"
              placeholder="1"
              after="NKN"
            />
          </ControlWrapper>
          <ControlWrapper title="Transfer fee">
            <template slot="descr">
              <div>
                <span> <b>Free:</b> {{ Number(0).toFixed(8) }} NKN</span><br />
                <span> <b>Low:</b> ~{{ lowFee.toFixed(8) }} NKN</span><br />
                <span> <b>Average:</b> ~{{ avgFee.toFixed(8) }} NKN</span><br />
                <span> <b>High:</b> ~{{ highFee.toFixed(8) }} NKN</span>
              </div>
            </template>
            <div class="settings__slider">
              <VueSlider
                class="settings__slider-control"
                v-if="showModal"
                v-model="transferFee"
                v-bind="sliderOptions"
              ></VueSlider>
            </div>
          </ControlWrapper>
        </template>
      </ModalGrid>
    </template>
    <template slot="footer">
      <Button theme="default" size="md" @click="closeModal">Close</Button>
      <Button
        v-if="activeTab === 'Transfers'"
        theme="default"
        size="md"
        @click="transfer(transferAddress, transferAmount, transferFee)"
        :disabled="transferDisabled"
        >Transfer</Button
      >
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
import Tabs from "@/components/Tabs/Tabs";
import Input from "@/components/Controls/Input/Input";

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
    Tabs,
    Input,
  },
  data: () => {
    return {
      tabs: ["General", "Transfers"],
      activeTab: "General",
      transferAddress: "",
      transferAmount: 0,
      walletAddress: "",
      walletBalance: 0,
      transferFee: 0,
      transferLoading: false,
      txFee: 0,
      avgFee: 0,
      lowFee: 0,
      highFee: 0,
      downloadPath: "",
      sliderOptions: {
        dotSize: 20,
        height: 2,
        tooltip: false,
        interval: 1,
        marks: {
          0: "Free",
          33: "Low",
          66: "Average",
          100: "High",
        },
        adsorb: true,
        included: true,
      },
    };
  },
  computed: {
    transferDisabled() {
      if (this.transferLoading || !this.isTransferAddress) {
        return true;
      } else {
        return false;
      }
    },

    isTransferAddress() {
      const address = this.transferAddress;
      const regexp = /^((^NKN([A-Za-z0-9]){33}){1})$/;
      return regexp.test(address);
    },
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
    openLink(url) {
      window.go.surge.MiddlewareFunctions.OpenLink(url);
    },
    transfer(address, amount, fee) {
      this.transferLoading = true;

      window.go.surge.MiddlewareFunctions.TransferToRecipient(
        address,
        amount.toString(),
        fee.toString()
      )
        .then((resp) => {
          console.log(resp);

          this.$notify({
            group: "notifications",
            text: `Transfer successful`,
            type: "success",
          });
        })
        .finally(() => {
          this.transferLoading = false;
        });
    },
    setMaxTransfer() {
      this.transferAmount = this.walletBalance;
    },
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
          this.$notify({
            group: "notifications",
            text: `Open API error: ` + err,
            type: "error",
          });
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
