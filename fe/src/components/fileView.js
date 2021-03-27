import { requestFileInfo } from 'api/storage'
import { getRootElement } from 'utils/element'

const createFileView = async (path) => {
  const { ok, status, content: fileInfo } = await requestFileInfo(path)
  if (!ok) {
    const tipContent = status === 404
      ? 'File not found'
      : 'Unexpected error'
    getRootElement()
      .removeAllChildren()
      .appendChild(document.createTextNode(tipContent))
    return
  }

  getRootElement()
    .removeAllChildren()
    .appendChild(document.createTextNode(`TODO: preview page for '${fileInfo.Name}'.`))
}

export { createFileView }
