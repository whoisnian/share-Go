const { htmlTemplatePlugin, copyPlugin } = require('./plugin')
const { fromRoot, fromOutput } = require('./function')
const { version } = require('../package.json')

const isProduction = process.env.NODE_ENV === 'production'

module.exports.buildConfig = {
  platform: 'browser',
  bundle: true,
  minify: isProduction,
  define: {
    __PACKAGE_VERSION__: `"${version}"`,
    __DEBUG__: `${!isProduction}`
  },
  entryPoints: [fromRoot('src/app.js')],
  entryNames: '[name]-[hash]',
  outdir: fromOutput('static'),
  logLevel: 'info',
  metafile: true,
  plugins: [
    htmlTemplatePlugin(fromOutput()),
    copyPlugin(fromRoot('public'), fromOutput())
  ]
}
