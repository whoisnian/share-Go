import { requestFileInfo } from 'api/storage'
import { createElement } from 'utils/element'
import { joinPath } from 'utils/function'

/** @param { string } oriPath */
const createFileView = async (oriPath) => {
  const { ok, status } = await requestFileInfo(oriPath)
  if (!ok) {
    return document.createTextNode(status === 404
      ? 'File not found'
      : 'Unexpected error'
    )
  }

  const iframe = createElement('iframe', {
    name: 'raw',
    src: joinPath('/api/raw', oriPath),
    style: 'position:fixed;border:none;top:0;bottom:0;left:0;right:0;width:100%;height:100%;'
  })

  const interval = setInterval(() => {
    if (window.frames['raw'] && window.frames['raw'].document.documentURI === iframe.src) {
      if (window.frames['raw'].document.body) window.frames['raw'].document.body.style = "margin:0;"
      clearInterval(interval)
    }
  }, 100)

  iframe.onload = () => {
    if (window.frames['raw'].document.body) window.frames['raw'].document.body.style = "margin:0;"
    clearInterval(interval)
  }

  return iframe
}

export { createFileView }
