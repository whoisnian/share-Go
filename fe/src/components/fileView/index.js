import { requestFileInfo } from 'api/storage'

const createFileView = async (path) => {
  const { ok, status, content: fileInfo } = await requestFileInfo(path)
  if (!ok) {
    return document.createTextNode(status === 404
      ? 'File not found'
      : 'Unexpected error'
    )
  }

  return document.createTextNode(`TODO: preview page for '${fileInfo.Name}'.`)
}

export { createFileView }
