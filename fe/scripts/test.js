import { readdirSync } from 'fs'
import { resolve } from 'path'
import { spawnSync } from 'child_process'
import { buildSync } from 'esbuild'
import { fromRoot } from './function.js'
import { buildConfig } from './esbuild.config.js'

const execNode = (code) => spawnSync('node', ['--enable-source-maps'], { cwd: fromRoot(), input: code })

const runTest = (path) => {
  if (path.endsWith('.test.js')) {
    const buildRes = buildSync({
      ...buildConfig(),
      platform: 'node',
      entryPoints: [path],
      outdir: undefined,
      outfile: undefined,
      plugins: undefined,
      write: false
      // sourcemap: 'inline'
    })
    if (buildRes.warnings.length > 0) console.error(buildRes.warnings)

    if (buildRes.outputFiles.length < 1) return
    const execRes = execNode(buildRes.outputFiles[0].text)
    console.log(execRes.stdout.toString())
    if (execRes.stderr.length > 0) {
      console.error(execRes.stderr.toString())
    }
  }
}

const walkDir = (dir, cb) => {
  const files = readdirSync(dir, { withFileTypes: true })
  files.forEach(file => {
    if (file.isDirectory()) walkDir(resolve(dir, file.name), cb)
    else cb(resolve(dir, file.name))
  })
}

const runMain = () => {
  walkDir(fromRoot('src'), runTest)
}

runMain()
