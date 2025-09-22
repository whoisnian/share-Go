import { createServer } from 'http'
import { context } from 'esbuild'
import httpProxy from 'http-proxy'
import { fromOutput } from './function.js'
import { buildConfig } from './esbuild.config.js'

const DEV_SERVER_PORT = 9100

// proxy_pass to target like nginx:
//   If the target is specified without a URI, the URI of the original request will be passed to the target.
//   If the target is specified with a URI, the request URI matching the pattern will be replaced before being passed to the target.
const routes = [
  // [ pattern, target ]
  ['^/api/', 'http://127.0.0.1:9000'],
  ['^/view/.*', '/static/index.html']
]

const createTransformer = (defaultTarget) => {
  const parsedRoutes = routes.map(([pattern, target]) => {
    const url = new URL(target, defaultTarget)
    const parsedTarget = url.origin
    const pathname = (target === '' || target === parsedTarget) ? null : url.pathname
    return [new RegExp(pattern), pathname, parsedTarget]
  })
  const tagDefault = '\x1b[32m[default]\x1b[0m'
  const tagProxy = '  \x1b[35m[proxy]\x1b[0m'
  const result = (req, target, originalUrl) => {
    console.log(`${target === defaultTarget ? tagDefault : tagProxy} ` +
      `${req.method} ${originalUrl}` +
      `${originalUrl !== req.url ? ` -> ${req.url}` : ''}` +
      `${target === defaultTarget ? '' : ` \x1b[36mto\x1b[0m ${target}`}`
    )
    return { req, target }
  }
  return (req) => {
    for (const [regexp, pathname, target] of parsedRoutes) {
      if (regexp.test(req.url)) {
        const originalUrl = req.url
        if (pathname) req.url = req.url.replace(regexp, pathname)
        return result(req, target, originalUrl)
      }
    }
    return result(req, defaultTarget, req.url)
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
  const targetTransform = createTransformer(`http://${hosts[0]}:${port}`)

  const proxy = httpProxy.createProxyServer({})
  const mainServer = createServer((req, res) => {
    if (req.url === '/') {
      res.writeHead(302, { location: '/view/' }).end()
      console.log(' \x1b[34m[direct]\x1b[0m redirect / to /view/')
      return
    }
    const { req: proxyReq, target } = targetTransform(req)
    proxy.web(proxyReq, res, { target })
  })
  mainServer.on('upgrade', (req, socket, head) => {
    const { req: proxyReq, target } = targetTransform(req)
    proxy.ws(proxyReq, socket, head, { target })
  })
  mainServer.listen(DEV_SERVER_PORT, () => {
    console.log('\x1b[32m>>>\x1b[0m The esbuild server has started serving build results.')
    console.log('\x1b[32m>>>\x1b[0m Starting the development server as a reverse proxy...\n')
    for (const host of hosts) {
      const matches = /(\d+)\.\d+\.\d+\.\d+/.exec(host)
      if (matches) {
        if (matches[1] === '127') console.log(` > Local:   http://${matches[0]}:${DEV_SERVER_PORT}`)
        else console.log(` > Network: http://${matches[0]}:${DEV_SERVER_PORT}`)
      }
    }
    console.log('')
  })
}

runMain()
