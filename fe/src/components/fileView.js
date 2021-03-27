import { requestFileInfo } from 'api/storage'
import { getRootElement, removeAllChildren } from 'utils/element'

const createFileView = async (path) => {
  const fileInfo = await requestFileInfo(path)

  const rootElement = getRootElement()
  removeAllChildren(rootElement)
  rootElement.appendChild(document.createTextNode(`TODO: preview page for '${fileInfo.Name}'.`))
}

export { createFileView }
