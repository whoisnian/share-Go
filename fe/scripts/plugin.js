const { writeFileSync, existsSync } = require('fs')
const { join, relative, extname } = require('path')
const { copyRecursivelySync } = require('./function')

/** @param {{ stylesheetList: string[], scriptList: string[] }} */
const generateHtmlFromTemplate = ({ stylesheetList, scriptList }) => {
  const stylesheetLink = stylesheetList.reduce((result, stylesheet) => {
    return result + `<link rel="stylesheet" href="${stylesheet}">`
  }, '')
  const scriptLink = scriptList.reduce((result, script) => {
    return result + `<script src="${script}"></script>`
  }, '')

  return `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>share-Go</title>
    ${stylesheetLink}
  </head>
  <body>
    <noscript>You need to enable JavaScript to run this app.</noscript>
    <main id="root"></main>
    ${scriptLink}
  </body>
</html>`
}

/** @param { string } base */
const htmlTemplatePlugin = (base) => ({
  name: 'htmlTemplatePlugin',
  setup (build) {
    const buildOptions = build.initialOptions
    const cwd = process.cwd()
    build.onEnd(result => {
      const options = {
        stylesheetList: [],
        scriptList: []
      }

      Object.keys(result.metafile.outputs).forEach((path) => {
        const absoluteLink = join('/', relative(base, join(cwd, path)))
        if (extname(path) === '.css') options.stylesheetList.push(absoluteLink)
        else if (extname(path) === '.js') options.scriptList.push(absoluteLink)
      })

      writeFileSync(join(buildOptions.outdir, 'index.html'), generateHtmlFromTemplate(options))
    })
  }
})

/**
 * @param { string } src
 * @param { string } dest
 */
const copyPlugin = (src, dest) => ({
  name: 'copyPlugin',
  setup (build) {
    build.onStart(() => {
      if (!src || !dest) throw new Error('SOURCE or DEST must not be blank')
      if (!existsSync(src)) throw new Error(`SOURCE '${src}' not exists`)

      copyRecursivelySync(src, dest)
    })
  }
})

module.exports = {
  generateHtmlFromTemplate,
  htmlTemplatePlugin,
  copyPlugin
}
