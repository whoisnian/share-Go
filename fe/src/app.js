import { init404Page } from 'pages/404'
import { initViewPage } from 'pages/view'
import { getPackageVersion } from 'utils/function'

__DEBUG__ && new EventSource('/esbuild').addEventListener('change', () => window.location.reload())

const routerList = [
  ['^/view(/|$)', initViewPage]
];

(() => {
  console.log(getPackageVersion())

  const pathname = window.location.pathname
  const [, handler] = routerList.find(([re]) => {
    return new RegExp(re).test(pathname)
  }) || [null, null]

  if (handler) handler()
  else init404Page()
})()
