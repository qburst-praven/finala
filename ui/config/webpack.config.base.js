const path = require("path");
const webpack = require("webpack");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const FaviconsWebpackPlugin = require("favicons-webpack-plugin");
const ESLintPlugin = require('eslint-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');

module.exports = (options) => ({
  mode: options.mode,
  entry: options.entry,
  devtool: options.devtool,
  output: Object.assign(
    {
      path: path.resolve(process.cwd(), "build", "static"),
      publicPath: "/static",
    },
    options.output
  ),
  module: {
    rules: options.module.rules.concat([
      {
        test: /\.js$/,
        exclude: /node_modules/,
        use: ["babel-loader"],
      },
      {
        test: /\.(css|sass|scss)$/,
        use: [
          MiniCssExtractPlugin.loader,
          "css-loader",
          "sass-loader",
        ],
      },
      {
        test: /\\.(png|jpe?g|gif)$/i,
        type: 'asset',
        parser: {
          dataUrlCondition: {
            maxSize: 8000,
          },
        },
        generator: {
          filename: 'images/[hash]-[name][ext]',
        },
      },
      {
        test: /\\.svg$/i,
        issuer: /\\.[jt]sx?$/,
        use: [{ loader: '@svgr/webpack', options: { icon: true, svgo: false } }],
      },
      {
        test: /\\.svg$/i,
        issuer: { not: /\\.[jt]sx?$/ },
        type: 'asset/resource',
        generator: {
          filename: 'images/[hash]-[name][ext]',
        },
      },
    ]),
  },
  plugins: options.plugins.concat([
    new ESLintPlugin({}),
    new MiniCssExtractPlugin({
      filename: '[name].[contenthash].css',
      chunkFilename: '[id].[contenthash].css',
    }),
    new FaviconsWebpackPlugin("./src/styles/icons/icon.png"),
    new CopyWebpackPlugin({
      patterns: [
        {
          from: path.resolve(process.cwd(), 'public'),
          to: path.resolve(process.cwd(), 'build'),
          globOptions: {
            ignore: ['**/index.html'],
          },
        },
      ],
    }),
  ]),
  resolveLoader: {
    modules: ["node_modules"],
  },
});
