let cssConfig = {
  loaderOptions: {
    sass: {
      data: `@import "@/assets/scss/_variables.scss";@import "@/assets/scss/_mixins.scss";`,
    },
  },
};

if (process.env.NODE_ENV == "production") {
  cssConfig.extract = {
    filename: "[name].css",
    chunkFilename: "[name].css",
  };
}

const path = require("path");
const PrerenderSPAPlugin = require("prerender-spa-plugin");

module.exports = {
  plugins: [
    new PrerenderSPAPlugin({
      // Required - The path to the webpack-outputted app to prerender.
      staticDir: path.join(__dirname, "dist"),
      // Required - Routes to render.
      routes: ["/", "/search", "/download", "/settings"],
    }),
  ],
};

module.exports = {
  chainWebpack: (config) => {
    let limit = 9999999999999999;
    const svgRule = config.module.rule("svg");

    svgRule.uses.clear();

    svgRule
      .test(/\.svg$/)
      .use("babel-loader")
      .loader("babel-loader")
      .end()
      .use("vue-svg-loader")
      .loader("vue-svg-loader");
    config.module
      .rule("images")
      .test(/\.(png|gif|jpg)(\?.*)?$/i)
      .use("url-loader")
      .loader("url-loader")
      .tap((options) => Object.assign(options, { limit: limit }));
    config.module
      .rule("fonts")
      .test(/\.(woff2?|eot|ttf|otf)(\?.*)?$/i)
      .use("url-loader")
      .loader("url-loader")
      .options({
        limit: limit,
      });
    config.module
      .rule("scss")
      .test(/\.scss$/)
      .use("vue-style-loader", "css-loader", "sass-loader")
      .loader("vue-style-loader", "css-loader", "sass-loader");
  },
  css: cssConfig,
  configureWebpack: {
    output: {
      filename: "[name].js",
    },
    optimization: {
      splitChunks: false,
    },
  },
  devServer: {
    disableHostCheck: true,
    host: "localhost",
  },
};
