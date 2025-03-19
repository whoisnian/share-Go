const { createServer, request: httpRequest } = require('http')
const { request: httpsRequest } = require('https')
const { context } = require('esbuild')
const { generateHtmlFromTemplate, copyPlugin } = require('./plugin')
const { fromRoot, fromOutput } = require('./function')
const { buildConfig } = require('./esbuild.config')

const request = (url, ...extraParams) => {
  return url.startsWith('https')
    ? httpsRequest(url, ...extraParams)
    : httpRequest(url, ...extraParams)
}

const proxyMap = [
  ['^/api', 'http://127.0.0.1:9000/api']
]

const proxyTransform = (url) => {
  if (!url) return null

  for (const [pattern, upstream] of proxyMap) {
    const re = new RegExp(pattern)
    if (re.test(url)) return url.replace(re, upstream)
  }
  return null
}

const runMain = async () => {
  const ctx = await context({
    ...buildConfig,
    write: false,
    sourcemap: true,
    entryNames: '[name]',
    plugins: [copyPlugin(fromRoot('public'), fromOutput())]
  })
  await ctx.watch()
  const { hosts, port } = await ctx.serve({ servedir: fromOutput() })

  createServer((req, res) => {
    if (req.url === '/') {
      res.writeHead(302, { location: '/view/' })
      res.end()
      console.log(' \x1b[34m[direct]\x1b[0m 302 redirect / to /view/')
      return
    }
    const proxyUrl = proxyTransform(req.url)
    const proxyReq = request(proxyUrl || `http://${hosts[0]}:${port}${req.url}`, {
      method: req.method,
      headers: req.headers
    }, proxyRes => {
      if (proxyRes.statusCode === 404 && proxyUrl === null) {
        res.writeHead(200, { 'content-type': 'text/html' })
        res.end(generateHtmlFromTemplate({
          stylesheetList: ['/static/app.css'],
          scriptList: ['/static/app.js']
        }))
        console.log(`\x1b[32m[esbuild]\x1b[0m 200 index ${req.method} ${req.url}`)
        return
      }

      res.writeHead(proxyRes.statusCode, proxyRes.headers)
      proxyRes.pipe(res, { end: true })
      if (proxyUrl === null) console.log(`\x1b[32m[esbuild]\x1b[0m ${res.statusCode} ok ${req.method} ${req.url}`)
      else console.log(`  \x1b[35m[proxy]\x1b[0m ${res.statusCode} ok ${req.method} ${proxyUrl}`)
    }).on('error', err => {
      res.writeHead(500, { 'content-type': 'text/plain' }).end(err.message)
      console.error(`  \x1b[31m[error]\x1b[0m 500 error ${req.method} ${proxyUrl || req.url} ${err.message}`)
    })

    req.pipe(proxyReq, { end: true })
  }).listen(9100, () => {
    console.log('esbuild dev server started: <http://127.0.0.1:9100>')
  })
}

runMain()
