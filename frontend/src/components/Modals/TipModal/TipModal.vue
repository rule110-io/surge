<template>
  <Modal :show.sync="showModal" @closeAndClear="closeAndClearModal">
    <template slot="title"> Tip seeders </template>
    <template slot="body">
      <div v-if="activeFile" class="modal__descr modal__item">
        You are going to tip <b>{{ activeFile.FileName }}</b> seeders. Please
        set the tip amount to be splitted among all the seeders
        <b>({{ activeFile.NumSeeders }})</b>.
      </div>
      <ModalGrid>
        <ControlWrapper title="Total tip amount">
          <Input
            v-model="amount"
            type="number"
            theme="light"
            size="md"
            placeholder="1"
            after="NKN"
          />
        </ControlWrapper>
        <ControlWrapper title="Transaction fee">
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
      </ModalGrid>
    </template>
    <template slot="footer">
      <Button theme="text" size="md" @click="closeAndClearModal">Close</Button>
      <Button theme="default" size="md" :disabled="amount <= 0" @click="tip"
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
  },
  watch: {
    showModal() {
      this.getAvgTxFee();
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
        .catch(() => {
          this.$store.dispatch("snackbar/updateSnack", {
            snack: `Open API error`,
            color: "error",
            timeout: false,
          });
        });
    },
    clearModal() {},
    tip() {
      console.log(this.activeFile);
      window.go.surge.MiddlewareFunctions.Tip(
        this.activeFile.FileHash,
        "" + this.amount,
        "" + this.txFee
      ).then(() => {
        this.closeModal();
        this.clearModal();
      });
    },
  },
};
</script>
