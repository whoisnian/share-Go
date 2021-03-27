const { resolve } = require('path')
const { version } = require('../package.json')

// const isDEV = process.env.NODE_ENV !== 'production'
const PATH_ROOT = resolve(__dirname, '..')
const PATH_OUTPUT = resolve(__dirname, '../dist')
const fromRoot = (...args) => resolve(PATH_ROOT, ...args)
const fromOutput = (...args) => resolve(PATH_OUTPUT, ...args)

module.exports.buildConfig = {
  platform: 'browser',
  bundle: true,
  define: { __PACKAGE_VERSION__: `"${version}"` },
  entryPoints: [fromRoot('src/app.js')],
  outfile: fromOutput('app.js')
}
