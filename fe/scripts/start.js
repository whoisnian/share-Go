const { createServer, request: httpRequest } = require('http')
const { request: httpsRequest } = require('https')
const { serve } = require('esbuild')
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

serve({
  servedir: fromOutput()
}, {
  ...buildConfig,
  entryNames: '[name]',
  plugins: [copyPlugin(fromRoot('public'), fromOutput())]
}).then(result => {
  createServer((req, res) => {
    if (req.url === '/') {
      res.writeHead(302, { 'location': '/view/' })
      res.end()
      return
    }
    const url = proxyTransform(req.url)
    const proxyReq = request(url || `http://${result.host}:${result.port}${req.url}`, {
      method: req.method,
      headers: req.headers
    }, proxyRes => {
      if (proxyRes.statusCode === 404 && url === null) {
        res.writeHead(200, { 'content-type': 'text/html' })
        res.end(generateHtmlFromTemplate({
          stylesheetList: ['/static/app.css'],
          scriptList: ['/static/app.js']
        }))
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
