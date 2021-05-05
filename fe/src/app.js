import { init404Page } from 'pages/404'
import { initViewPage } from 'pages/view'
import { getRootElement } from 'utils/element'
import { getPackageVersion } from 'utils/function'
import { observer } from 'utils/observer'

const routerList = [
  ['^/view(/|$)', initViewPage]
];

(() => {
  console.log(getPackageVersion())

  window.observer = observer

  const changePage = (_, pathname) => {
    const [, handler] = routerList.find(([re]) => {
      return new RegExp(re).test(pathname)
    }) || [null, null]
    observer.set('page', handler || init404Page)
  }
  observer.onchange('pathname', changePage)

  const render = async (_, initPage) => {
    getRootElement()
      .removeAllChildren()
      .appendChild(await initPage())
  }
  observer.onchange('page', render)

  observer.set('pathname', window.location.pathname)
})()
