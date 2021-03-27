import { requestListDir } from 'api/storage'
import { getRootElement, removeAllChildren } from 'utils/element'

const createDirView = async (path) => {
  const fileInfos = await requestListDir(path)
  fileInfos.sort((a, b) => {
    if (a.Type === b.Type) return a.Name.localeCompare(b.Name, 'zh-CN')
    return a.Type < b.Type
  })

  const rootElement = getRootElement()
  removeAllChildren(rootElement)
  fileInfos.forEach(({ Type, Name, Size, MTime }) => {
    const item = document.createElement('li')
    item.innerHTML = `${Type} ${Name} ${Size} ${MTime}`
    rootElement.appendChild(item)
  })
}

export { createDirView }
