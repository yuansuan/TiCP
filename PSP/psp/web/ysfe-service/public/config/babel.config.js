module.exports = {
  babelrc: false,
  configFile: false,
  presets: [require.resolve('babel-preset-react-app')],
  plugins: [
    [
      require.resolve('babel-plugin-named-asset-import'),
      {
        loaderMap: {
          svg: {
            ReactComponent: '@svgr/webpack?-svgo,+titleProp,+ref![path]',
          },
        },
      },
    ],
    'react-hot-loader/babel',
    [
      require.resolve('babel-plugin-import'),
      {
        libraryName: 'antd',
        libraryDirectory: 'es',
        style: true,
      },
      'antd',
    ],
    [
      require.resolve('babel-plugin-import'),
      {
        libraryName: '@ys/components',
        camel2DashComponentName: false,
        customName: (name) => `@ys/components/dist/${name}`,
      },
      '@ys/components',
    ],
    require.resolve('@babel/plugin-syntax-dynamic-import'),
  ],
}
