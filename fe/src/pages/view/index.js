import { createFileView } from 'components/fileView'
import { createDirView } from 'components/dirView'
import { FileType, requestFileInfo } from 'api/storage'
import { getRootElement, removeAllChildren } from 'utils/element'

const initViewPage = async () => {
  const pathname = window.location.pathname
  const path = pathname.slice(pathname.indexOf('/view') + 5)

  const { ok, status, content: fileInfo } = await requestFileInfo(path)
  if (ok && fileInfo.Type === FileType.typeRegular) {
    await createFileView(path)
  } else if (ok && fileInfo.Type === FileType.typeDirectory) {
    await createDirView(path)
  } else {
    const tipContent = ok
      ? 'Can not parse fileInfo'
      : status === 404
        ? 'File or dir not found'
        : 'Unexpected error'
    const rootElement = getRootElement()
    removeAllChildren(rootElement)
    rootElement.appendChild(document.createTextNode(tipContent))
  }
}

export { initViewPage }
