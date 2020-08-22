<template>
  <v-tour
    name="myTour"
    :steps="steps"
    :options="options"
    :callbacks="callbacks"
  ></v-tour>
</template>

<style lang="scss">
@import "./Tour.scss";
</style>

<script>
export default {
  components: {},
  data: () => ({
    callbacks: {},
    options: {
      useKeyboardNavigation: false,
      labels: {
        buttonSkip: "Skip tour",
        buttonPrevious: "Previous",
        buttonNext: "Next",
        buttonStop: "Finish",
      },
    },
    steps: [
      {
        target: "#app",
        header: {
          title: "Hello there, welcome to surge!",
        },
        content: `It seems like you have opened the app for the first time – Would you like to show you around a little?`,
        params: {
          highlight: false,
          enableScrolling: false,
        },
      },
      {
        target: ".sidebar",
        header: {
          title: "Navigation panel",
        },
        content: `Let’s get over the basics. On the left is the navigation panel. It lets you switch between the “Explore”, “Files” and “Settings” page`,
        params: {
          highlight: true,
          enableScrolling: false,
        },
      },
      {
        target: ".header",
        header: {
          title: "Top panel",
        },
        content: `On the top panel you can easily search for files published in the NKN blockchain. Also there are quick-actions for switching the theme between dark and night-mode and accessing the latest events`,
        params: {
          highlight: true,
          enableScrolling: false,
        },
      },
      {
        target: ".network-stats",
        header: {
          title: "Bottom panel",
        },
        content: `On the bottom panel you may add your local files so that they can be shared with your friends and family. In addition there is a graph that shows you the download and upload speed of your client.`,
        params: {
          highlight: true,
          enableScrolling: false,
        },
      },
      {
        target: ".table",
        header: {
          title: "That’s for the basics",
        },
        content: `Should we also go through each single page so I can show you all the features?`,
        params: {
          highlight: false,
          enableScrolling: false,
        },
      },
    ],
  }),
  mounted() {
    this.init();
  },
  methods: {
    init() {
      this.callbacks.onStop = () => {
        this.$store.dispatch("tour/offTour");
      };
      this.callbacks.onNextStep = (currentStep) => {
        if (currentStep === 2) {
          this.$router.push("download");
        }
      };

      this.$tours["myTour"].start();
    },
  },
};
</script>
