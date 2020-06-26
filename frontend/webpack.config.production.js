const webpack = require('webpack');
const path = require('path');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const UglifyJsPlugin = require('uglifyjs-webpack-plugin');
const OptimizeCSSAssetsPlugin = require("optimize-css-assets-webpack-plugin");
const CleanWebpackPlugin = require('clean-webpack-plugin');
const OptimizeCssAssetsPlugin = require('optimize-css-assets-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const WebpackShellPlugin = require("webpack-shell-plugin");

module.exports = {
	target: "web",
	mode: 'production',
	devtool: 'source-map',

	//entry: path.join(__dirname, '/src/app/index.js'),
	entry: path.resolve('./src/app/index.js'),
	output: {
		//path: path.resolve(__dirname, 'dist'),
		path: path.resolve(__dirname, 'dist'),
		filename: 'main.bundle.js'
	},

	module: {
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
				test: /\.(jpe?g|png|gif|ttf|woff|eot|woff2|svg)$/i,
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
	]},

	resolve: {
		alias: {
			jquery: "jquery/dist/jquery",
			index: path.resolve(__dirname, 'src/app/index.js')
		}
	},

	optimization: {
		minimizer: [
			new UglifyJsPlugin({
				cache: true,
				parallel: true,
				sourceMap: false
			}),
			new OptimizeCSSAssetsPlugin({})
		]
	},
  	plugins: [
		new CleanWebpackPlugin(['../dist']),
		new webpack.ProvidePlugin({
			$: 'jquery',
			jQuery: 'jquery',
		}),

		new MiniCssExtractPlugin({
			filename: 'main.css',
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
