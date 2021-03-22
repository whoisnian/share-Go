const { build } = require('esbuild')
const { buildConfig } = require('./esbuild.config')

build(buildConfig).catch(() => process.exit(1))
