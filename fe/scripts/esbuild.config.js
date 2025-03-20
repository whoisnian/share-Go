import { htmlTemplatePlugin, copyPlugin } from './plugin.js'
import { fromRoot, fromOutput } from './function.js'
import packageJson from '../package.json' with { type: 'json' }

const { version } = packageJson
const isProduction = process.env.NODE_ENV === 'production'

const buildConfig = () => ({
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
})

export {
  buildConfig
}
