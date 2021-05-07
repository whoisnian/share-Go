const { createServer, request: httpRequest } = require('http')
const { request: httpsRequest } = require('https')
const { resolve } = require('path')
const { serve } = require('esbuild')
const { buildConfig } = require('./esbuild.config')

const PATH_ROOT = resolve(__dirname, '..')
// const PATH_OUTPUT = resolve(__dirname, '../dist')
const fromRoot = (...args) => resolve(PATH_ROOT, ...args)
// const fromOutput = (...args) => resolve(PATH_OUTPUT, ...args)

const exampleHtml = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Example</title>
    <link rel="stylesheet" href="/app.css">
  </head>
  <body>
    <noscript>You need to enable JavaScript to run this app.</noscript>
    <main id="root"></main>
    <script src="/app.js"></script>
  </body>
</html>`

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

serve({
  servedir: fromRoot('public')
}, {
  ...buildConfig,
  outfile: fromRoot('public/app.js')
}).then(result => {
  createServer((req, res) => {
    const url = proxyTransform(req.url)
    const proxyReq = request(url || `http://${result.host}:${result.port}${req.url}`, {
      method: req.method,
      headers: req.headers
    }, proxyRes => {
      if (proxyRes.statusCode === 404 && url === null) {
        res.writeHead(200, { 'content-type': 'text/html' })
        res.end(exampleHtml)
        return
      }

      res.writeHead(proxyRes.statusCode, proxyRes.headers)
      proxyRes.pipe(res, { end: true })
    })

    req.pipe(proxyReq, { end: true })
  }).listen(9100, () => {
    console.log('esbuild dev server started: <http://127.0.0.1:9100>')
  })
})
