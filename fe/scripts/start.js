import { createServer } from 'http'
import { context } from 'esbuild'
import httpProxy from 'http-proxy'
import { fromOutput } from './function.js'
import { buildConfig } from './esbuild.config.js'

const routes = [
  ['^/api/', 'http://127.0.0.1:9000'],
  ['^/view/.*', '/static/index.html']
]

const createTransformer = (esbuildUpstream) => {
  const parsedRoutes = routes.map(([pattern, upstream]) => {
    const url = new URL(upstream, esbuildUpstream)
    const target = url.origin
    const path = (upstream === '' || upstream === target) ? null : url.pathname
    return [pattern, path, target]
  })
  return (req) => {
    for (const [pattern, path, target] of parsedRoutes) {
      const regexp = new RegExp(pattern)
      if (regexp.test(req.url)) {
        if (path) req.url = req.url.replace(regexp, path)
        return { req, target }
      }
    }
    return { req, target: esbuildUpstream }
  }
}

const runMain = async () => {
  const ctx = await context({
    ...buildConfig(),
    write: false,
    sourcemap: true,
    entryNames: '[name]'
  })
  await ctx.watch()
  const { hosts, port } = await ctx.serve({ servedir: fromOutput() })
  const esbuildUpstream = `http://${hosts[0]}:${port}`
  const targetTransform = createTransformer(esbuildUpstream)

  const proxy = httpProxy.createProxyServer({})
  const mainServer = createServer((req, res) => {
    const originalUrl = req.url
    if (originalUrl === '/') {
      res.writeHead(302, { location: '/view/' }).end()
      console.log(' \x1b[34m[direct]\x1b[0m redirect / to /view/')
      return
    }

    const { req: proxyReq, target } = targetTransform(req)
    if (target === esbuildUpstream) console.log(`\x1b[32m[esbuild]\x1b[0m ${req.method} ${originalUrl}${originalUrl !== proxyReq.url ? ` -> ${proxyReq.url}` : ''}`)
    else console.log(`  \x1b[35m[proxy]\x1b[0m ${req.method} ${originalUrl}${originalUrl !== proxyReq.url ? ` -> ${proxyReq.url}` : ''} to ${target}`)
    proxy.web(proxyReq, res, { target })
  })
  mainServer.on('upgrade', (req, socket, head) => {
    const originalUrl = req.url
    const { req: proxyReq, target } = targetTransform(req)
    if (target === esbuildUpstream) console.log(`\x1b[32m[esbuild]\x1b[0m ${req.method} ${originalUrl}${originalUrl !== proxyReq.url ? ` -> ${proxyReq.url}` : ''}`)
    else console.log(`  \x1b[35m[proxy]\x1b[0m ${req.method} ${originalUrl}${originalUrl !== proxyReq.url ? ` -> ${proxyReq.url}` : ''} to ${target}`)
    proxy.ws(proxyReq, socket, head, { target })
  })
  mainServer.listen(9100, () => {
    console.log('esbuild dev server started: <http://127.0.0.1:9100>')
  })
}

runMain()
