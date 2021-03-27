import { requestListDir } from 'api/storage'
import { getRootElement } from 'utils/element'

const createDirView = async (path) => {
  const { ok, status, content } = await requestListDir(path)
  if (!ok) {
    const tipContent = status === 404
      ? 'Dir not found'
      : 'Unexpected error'
    getRootElement()
      .removeAllChildren()
      .appendChild(document.createTextNode(tipContent))
    return
  }

  const fileInfos = content.FileInfos
  fileInfos.sort((a, b) => {
    if (a.Type === b.Type) return a.Name.localeCompare(b.Name, 'zh-CN')
    return a.Type < b.Type
  })

  const rootElement = getRootElement().removeAllChildren()
  fileInfos.forEach(({ Type, Name, Size, MTime }) => {
    const item = document.createElement('li')
    item.innerHTML = `${Type} ${Name} ${Size} ${MTime}`
    rootElement.appendChild(item)
  })
}

export { createDirView }
