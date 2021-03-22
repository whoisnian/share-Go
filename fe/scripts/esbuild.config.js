const { resolve } = require('path')

const isDEV = process.env.NODE_ENV !== 'production'
const PATH_ROOT = resolve(__dirname, '..')
const PATH_OUTPUT = resolve(__dirname, '../dist')
const fromRoot = (...args) => resolve(PATH_ROOT, ...args)
const fromOutput = (...args) => resolve(PATH_OUTPUT, ...args)

module.exports.buildConfig = {
  platform: 'browser',
  entryPoints: [fromRoot('src/app.js')],
  bundle: isDEV,
  outfile: fromOutput('bundle.js')
}
