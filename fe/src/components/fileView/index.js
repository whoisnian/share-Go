import { requestFileInfo } from 'api/storage'

/** @param { string } oriPath */
const createFileView = async (oriPath) => {
  const { ok, status, content: fileInfo } = await requestFileInfo(oriPath)
  if (!ok) {
    return document.createTextNode(status === 404
      ? 'File not found'
      : 'Unexpected error'
    )
  }

  return document.createTextNode(`TODO: preview page for '${fileInfo.Name}'.`)
}

export { createFileView }
