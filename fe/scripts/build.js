const { build } = require('esbuild')
const { buildConfig } = require('./esbuild.config')

build(buildConfig).catch((err) => {
  console.error(err)
  process.exit(1)
})
