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
        content: `Would you like to show you around a little?`,
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
        target: "#app",
        header: {
          title: "That’s for the basics",
        },
        content: `Should we also go through each single page so I can show you all the features?`,
        params: {
          highlight: false,
          enableScrolling: false,
        },
      },
      {
        target: "#search",
        header: {
          title: 'Let\'s check out the "Search" view',
        },
        content: `It makes it easy for you to search the NKN network for files shared by others through surge.`,
        params: {
          highlight: true,
          enableScrolling: false,
        },
      },
      {
        target: "#search_results",
        header: {
          title: "On top you see the search results",
        },
        content: `Here you can see all the files that are currently shared by any user in the NKN network. This table will auto update when new files are found.`,
        params: {
          highlight: true,
          enableScrolling: false,
        },
      },
      {
        target: "#search_input",
        header: {
          title: "Dont forget!",
        },
        content: `You can always use the search bar to filter the results if you're searching for something specific.`,
        params: {
          highlight: true,
          enableScrolling: false,
        },
      },
      {
        target: "#download",
        header: {
          title: 'Now lets have a look at the "Downloads" page',
        },
        content: `You can see the status of each file you interact with here.`,
        params: {
          highlight: true,
          enableScrolling: false,
        },
      },
      {
        target: "#files_table",
        header: {
          title: "This is the files table",
        },
        content: `Here you can see the status of all your local files. The status could vary by "Downloading", "Seeding" or "Finished".`,
        params: {
          highlight: true,
          enableScrolling: false,
        },
      },
      {
        target: "#settings",
        header: {
          title: 'Last thing on the list is the "Settings" page',
        },
        content: `Since this page changes so often you're free to explore it by yourself.`,
        params: {
          highlight: true,
          enableScrolling: false,
        },
      },
      {
        target: "#app",
        header: {
          title: "That's all for now!",
        },
        content: `We hope you now have a better understanding on how things work. Feek free to join our discord when you got problems or want to participate in the project!`,
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
        if (currentStep === 4) {
          this.$router.push("search");
        }
        if (currentStep === 8) {
          this.$router.push("download");
        }
        if (currentStep === 10) {
          this.$router.push("settings");
        }
      };

      this.$tours["myTour"].start();
    },
  },
};
</script>
