import { createFileView } from 'components/fileView'
import { createDirView } from 'components/dirView'
import { FileType, requestFileInfo } from 'api/storage'
import { getRootElement } from 'utils/element'

const initViewPage = async () => {
  const pathname = window.location.pathname
  const path = pathname.slice(pathname.indexOf('/view') + 5)

  let mainElement
  const { ok, status, content: fileInfo } = await requestFileInfo(path)
  if (ok && fileInfo.Type === FileType.typeRegular) {
    mainElement = await createFileView(path)
  } else if (ok && fileInfo.Type === FileType.typeDirectory) {
    mainElement = await createDirView(path)
  } else {
    mainElement = document.createTextNode(ok
      ? 'Can not parse fileInfo'
      : status === 404
        ? 'File or dir not found'
        : 'Unexpected error'
    )
  }

  getRootElement()
    .removeAllChildren()
    .appendChild(mainElement)
}

export { initViewPage }
