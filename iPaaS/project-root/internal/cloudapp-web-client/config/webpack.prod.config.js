const HtmlWebpackPlugin = require('html-webpack-plugin')
const CopyWebpackPlugin = require('copy-webpack-plugin')

const path = require('path')

const currentDate = new Date()
const year = currentDate.getFullYear()
const month = String(currentDate.getMonth() + 1).padStart(2, '0')
const day = String(currentDate.getDate()).padStart(2, '0')

module.exports = {
  mode: 'production',
  entry: './app/app.js',
  output: {
    filename: `lightdesk.${year}-${month}-${day}.[hash].js`,
    path: path.resolve(__dirname, '../dist')
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, '../app/')
    }
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: path.join(__dirname, '../app/tmpl/index.html'),
      filename: 'index.html',
      favicon: path.join(__dirname, '../resource/images/favicon.ico')
    }),
    new CopyWebpackPlugin([
      {
        from: path.join(__dirname, '../resource/images'),
        to: path.join(__dirname, '../dist/static'),
        toType: 'dir'
      }
    ])
  ],
  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
          options: {
            presets: ['@babel/preset-env']
          }
        }
      },
      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader']
      },
      {
        test: /\.scss$/,
        use: ['style-loader', 'css-loader', 'sass-loader']
      }
    ]
  },
  devServer: {
    https: true,
    open: true
  }
}
