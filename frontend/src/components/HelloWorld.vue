<template>
  <div class="container">
    <h1>{{message}}</h1>
    <a @click="getMessage">Press Me!</a>
  </div>
</template>

<style lang="scss">
@import './HelloWorld.scss';
</style>

<script>
export default {
  data() {
    return {
      message: " "
    };
  },
  methods: {
    getMessage: function() {
      var self = this;
      window.backend.basic().then(result => {
        self.message = result;
      });
    }
  }
};
</script>