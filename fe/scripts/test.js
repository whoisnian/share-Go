const { readdirSync } = require('fs')
const { resolve } = require('path')
const { spawnSync } = require('child_process')
const { buildSync } = require('esbuild')
const { buildConfig } = require('./esbuild.config')

const PATH_ROOT = resolve(__dirname, '..')
const fromRoot = (...args) => resolve(PATH_ROOT, ...args)
const execNode = (code) => spawnSync('node', ['--enable-source-maps'], { cwd: fromRoot(), input: code })

const runTest = (path) => {
  if (path.endsWith('.test.js')) {
    const buildRes = buildSync({
      ...buildConfig,
      platform: 'node',
      entryPoints: [path],
      outfile: undefined,
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
