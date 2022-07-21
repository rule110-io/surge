<template>
  <Modal :show.sync="showModal" @closeAndClear="closeAndClearModal">
    <template slot="title"> Tip seeders </template>
    <template slot="body">
      <div v-if="activeFile" class="modal__descr modal__item">
        You are going to tip <b>{{ activeFile.FileName }}</b> seeders. Please
        set the tip amount to be split among all the seeders
        <b>({{ activeFile.NumSeeders }})</b>.
      </div>
      <ModalGrid>
        <ControlWrapper title="Total tip amount">
          <template slot="descr">
            <div>Balance: {{ walletBalance.toFixed(8) }}</div>
          </template>
          <Input
            v-model="amount"
            type="number"
            theme="light"
            size="md"
            placeholder="1"
            after="NKN"
          />
        </ControlWrapper>
        <ControlWrapper
          :title="`Transaction fee (for ${activeFile.NumSeeders} seeders)`"
        >
          <template slot="descr">
            <div>
              <span>
                <b>Free:</b>
                {{ Number(0).toFixed(8) }} NKN</span
              ><br />
              <span>
                <b>Low:</b> ~{{
                  (lowFee * activeFile.NumSeeders).toFixed(8)
                }}
                NKN</span
              ><br />
              <span>
                <b>Average:</b> ~{{
                  (avgFee * activeFile.NumSeeders).toFixed(8)
                }}
                NKN</span
              ><br />
              <span>
                <b>High:</b> ~{{
                  (highFee * activeFile.NumSeeders).toFixed(8)
                }}
                NKN</span
              >
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
      <Button theme="text" size="md" @click="closeAndClearModal">Close</Button>
      <Button theme="default" size="md" :disabled="disabled" @click="tip"
        >Tip seeders</Button
      >
    </template>
  </Modal>
</template>

<script>
import { mapState } from "vuex";

import FormMixin from "@/mixins/FormMixin.js";

import Modal from "@/components/Modals/Modal/Modal";
import Button from "@/components/Button/Button";
import ModalGrid from "@/components/Modals/ModalGrid/ModalGrid";
import ControlWrapper from "@/components/Controls/ControlWrapper/ControlWrapper";
import Input from "@/components/Controls/Input/Input";

import VueSlider from "vue-slider-component";

import axios from "axios";

export default {
  mixins: [FormMixin],
  components: { Modal, Button, ModalGrid, ControlWrapper, Input, VueSlider },
  props: {
    file: {
      type: Object,
      default: () => {},
    },
  },
  data: () => {
    return {
      loading: false,
      walletBalance: 0,
      amount: 1,
      txFee: 66,
      avgFee: 0,
      lowFee: 0,
      highFee: 0,
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
    ...mapState("files", ["activeFile", "localFilesConfig"]),
    selectedFee() {
      let fee = 0.1;

      switch (this.txFee) {
        case 0:
          fee = this.lowFee;
          break;
        case 50:
          fee = this.avgFee;
          break;
        case 100:
          fee = this.highFee;
          break;
      }

      return fee;
    },
    disabled() {
      if (this.amount <= 0 || this.loading) {
        return true;
      } else {
        return false;
      }
    },
  },
  watch: {
    showModal() {
      this.getAvgTxFee();
      this.getWalletBalance();
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
          this.$notify({
            group: "notifications",
            text: `Open API error: ` + err,
            type: "error",
          });
        });
    },
    getWalletBalance() {
      window.go.surge.MiddlewareFunctions.GetWalletBalance().then((resp) => {
        this.walletBalance = parseFloat(resp);
      });
    },
    clearModal() {},
    tip() {
      this.loading = true;

      window.go.surge.MiddlewareFunctions.Tip(
        this.activeFile.FileHash,
        this.amount.toString(),
        this.txFee.toString()
      )
        .then(() => {
          this.closeModal();
          this.clearModal();

          this.$notify({
            group: "notifications",
            text: `Tip successful`,
            type: "success",
          });
        })
        .finally(() => {
          this.loading = false;
        });
    },
  },
};
</script>
