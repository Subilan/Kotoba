const path = require("path");
const webpack = require("webpack");

module.exports = {
  entry: "./src/dist.ts",
  mode: "production",
  target: "browserslist",
  module: {
    rules: [
      {
        test: /\.ts$/,
        use: "babel-loader",
        exclude: /node_modules/
      },
    ],
  },
  resolve: {
    extensions: [".ts", ".js"],
    fallback: {
        "buffer": require.resolve("buffer")
    }
  },
  output: {
    filename: "kotoba.js",
    path: path.resolve(__dirname, "dist"),
  },
  plugins: [
    new webpack.ProvidePlugin({
        Buffer: ['buffer', 'Buffer']
    }),
    new webpack.ProvidePlugin({
        process: 'process/browser'
    })
  ]
};
