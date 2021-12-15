<template>
  <TransfersInfoCard
    class="transfers-details"
    :active="fileDetails"
    title="Details"
  >
    <template slot="info">
      <div class="transfers-details__tabs">
        <div
          class="transfers-details__item"
          @click="setActiveTab('General')"
          :class="{ 'transfers-details__item_active': activeTab === 'General' }"
        >
          General
        </div>
        <div
          class="transfers-details__item"
          @click="setActiveTab('Seeders')"
          :class="{ 'transfers-details__item_active': activeTab === 'Seeders' }"
        >
          Seeders
        </div>
      </div>
    </template>

    <template slot="body" v-if="activeFileDetails">
      <div v-show="activeTab === 'General'">
        <FileChunks
          class="transfers-details__chunks"
          :file="activeFileDetails"
        />
        <div class="transfers-details__info">
          <div class="transfers-details__info-item">
            <div class="transfers-details__info-left">Downloaded</div>
            <div class="transfers-details__info-right">
              {{ activeFileDetails.BytesDownloaded | prettyBytes(0) }}
            </div>
          </div>
          <div class="transfers-details__info-item">
            <div class="transfers-details__info-left">Uploaded</div>
            <div class="transfers-details__info-right">
              {{ activeFileDetails.BytesUploaded | prettyBytes(0) }}
            </div>
          </div>
          <div class="transfers-details__info-item">
            <div class="transfers-details__info-left">Date</div>
            <div class="transfers-details__info-right">
              {{
                $moment(activeFileDetails.DateTimeAdded * 1000).format(
                  "DD.MM.YYYY"
                )
              }}
            </div>
          </div>
          <div class="transfers-details__info-item">
            <div class="transfers-details__info-left">Topic</div>
            <div class="transfers-details__info-right">
              #{{ activeFileDetails.Topic }}
            </div>
          </div>
          <div class="transfers-details__info-item">
            <div class="transfers-details__info-left">Total Chunks</div>
            <div class="transfers-details__info-right">
              {{ activeFileDetails.NumChunks }}
            </div>
          </div>
          <div class="transfers-details__info-item">
            <div class="transfers-details__info-left">Chunks Downloaded</div>
            <div class="transfers-details__info-right">
              {{ activeFileDetails.ChunksDownloaded }}
            </div>
          </div>
          <div class="transfers-details__info-item">
            <div class="transfers-details__info-left">Chunks Shared</div>
            <div class="transfers-details__info-right">
              {{ activeFileDetails.ChunksShared }}
            </div>
          </div>
          <div class="transfers-details__info-item">
            <div class="transfers-details__info-left">Progress</div>
            <div class="transfers-details__info-right">
              {{ (activeFileDetails.Progress * 100).toFixed(2) }}%
            </div>
          </div>
        </div>
      </div>
      <div class="transfers-details__seeders" v-show="activeTab === 'Seeders'">
        <div class="transfers-details__seeders-row">
          <div class="transfers-details__seeders-heading">#</div>
          <div class="transfers-details__seeders-heading">Pubkey</div>
        </div>
        <div
          class="transfers-details__seeders-row"
          v-for="(item, i) in activeFileDetails.Seeders"
          :key="i"
        >
          <div class="transfers-details__seeders-value">
            {{ i + 1 }}
          </div>
          <div class="transfers-details__seeders-value">
            {{ item }}
          </div>
        </div>
      </div>
    </template>
  </TransfersInfoCard>
</template>

<style lang="scss">
@import "./TransfersDetails";
</style>

<script>
import { mapState } from "vuex";

import TransfersInfoCard from "@/components/Widgets/Transfers/TransfersInfoCard/TransfersInfoCard";
import FileChunks from "@/components/File/FileChunks/FileChunks";

export default {
  components: { TransfersInfoCard, FileChunks },
  computed: {
    ...mapState("files", ["fileDetails", "selectedFiles"]),
    ...mapState("globalBandwidth", ["statusBundle"]),
  },
  data: () => {
    return {
      activeTab: "General",
      activeFileDetails: null,
    };
  },

  watch: {
    selectedFiles(newItems) {
      if (!newItems.length) {
        this.lastSelected = null;
        return;
      }

      const lastSelected = newItems[newItems.length - 1];

      this.getActiveFileDetails(lastSelected);
    },
    statusBundle(newEvent) {
      if (!this.activeFileDetails) return;

      const { FileHash } = this.activeFileDetails;
      const newFileHash = this._.find(newEvent, { FileHash });
      const isNewFileHash = !this._.isEmpty(newFileHash);

      if (isNewFileHash) {
        this.getActiveFileDetails({
          ...this.activeFileDetails,
          ...newFileHash,
        });
      }
    },
  },
  mounted() {},
  methods: {
    getActiveFileDetails(file) {
      window.go.surge.MiddlewareFunctions.GetFileDetails(file.FileHash).then(
        (resp) => {
          this.activeFileDetails = { ...file, ...resp };

          console.log(this.activeFileDetails);
        }
      );
    },
    setActiveTab(str) {
      this.activeTab = str;
    },
  },
};
</script>