import { build } from 'esbuild'
import { buildConfig } from './esbuild.config.js'

build(buildConfig()).catch((err) => {
  console.error(err)
  process.exit(1)
})
