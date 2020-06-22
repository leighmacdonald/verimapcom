const webpack = require('webpack');
const path = require('path');
const BundleAnalyzerPlugin = require('webpack-bundle-analyzer').BundleAnalyzerPlugin;
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const OptimizeCssAssetsPlugin = require('optimize-css-assets-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const WebpackShellPlugin = require('webpack-shell-plugin');

module.exports = {
    target: "web",
    mode: 'development',

    entry: path.join(__dirname, '/src/app/index.js'),

    output: {
        path: path.resolve(__dirname, 'dist'),
        filename: 'main.bundle.js'
    },

    devServer: {
        port: 9000, // port to run dev-server
        contentBase: path.join(__dirname, 'src/public'),
        watchContentBase: true,
        watchOptions: {
            poll: true
        }
    },

    devtool: 'source-map',

    module: {  // where we defined file patterns and their loaders
        rules: [
            {
                test: /\.js$/,
                exclude: /(node_modules|bower_components)/,
                use: [{
                    loader: 'babel-loader',
                    options: {
                        presets: ['@babel/preset-env']
                    }
                }
                ]
            },
            {
                test: /\.(sa|sc|c)ss$/,
                use: [
                    MiniCssExtractPlugin.loader,
                    'css-loader',
                    'postcss-loader',
                    'sass-loader',
                ],
            },
            {
                test: /\.(jpg|png|gif|svg)$/,
                loader: 'image-webpack-loader',
                // Specify enforce: 'pre' to apply the loader
                // before url-loader/svg-url-loader
                // and not duplicate it in rules with them
                enforce: 'pre'
            },
            {
                test: /\.(jpe?g|png|gif)$/i,
                use: [
                    {
                        loader: 'url-loader',
                        options: {
                            limit: 10000,
                            name: "images/[name].[ext]"
                        },
                    }
                ]
            }
        ]
    },

    optimization: {
        minimize: false
    },
    plugins: [
        new webpack.ProvidePlugin({
            $: 'jquery',
            jQuery: 'jquery',
        }),
        new MiniCssExtractPlugin({
            filename: '[name].css',
            chunkFilename: '[id].css'
        }),
        new OptimizeCssAssetsPlugin({
            cssProcessorPluginOptions: {
                preset: ['default', { discardComments: { removeAll: true } }]
            }
        }),
        new CopyWebpackPlugin([
            {
                from: 'src/public/fonts',
                to: 'fonts'
            }
        ]),
        new WebpackShellPlugin({
            //onBuildStart:['echo "Webpack Start"'],
            onBuildEnd:['bash compress.sh']})
    ]

};
