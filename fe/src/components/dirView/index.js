import { requestListDir } from 'api/storage'
import { createIcon } from 'components/icon'
import { createElement } from 'utils/element'
import { calcFromBytes, calcRelativeTime } from 'utils/function'
import './style.css'

/** @param { string } path */
const createDirView = async (path) => {
  const { ok, status, content } = await requestListDir(path)
  if (!ok) {
    return document.createTextNode(status === 404
      ? 'Dir not found'
      : 'Unexpected error'
    )
  }

  const fileInfos = content.FileInfos
  fileInfos.sort((a, b) => {
    if (a.Type === b.Type) return a.Name.localeCompare(b.Name, 'zh-CN')
    return a.Type < b.Type
  })

  const mainElement = createElement('div', { class: 'DirView-main' })
  const header = createElement('div', { class: 'DirView-header' })
  const currentPath = createElement('div', { class: 'DirView-currentPath' })
  currentPath.innerHTML = decodeURIComponent(path)

  const parentIcon = createIcon('folder-parent')
  parentIcon.onclick = () => { window.history.back() }
  header.appendChild(parentIcon)
  header.appendChild(createIcon('refresh'))
  header.appendChild(createIcon('home'))
  header.appendChild(currentPath)
  header.appendChild(createIcon('paste'))
  header.appendChild(createIcon('folder-new'))
  header.appendChild(createIcon('file-new'))
  header.appendChild(createIcon('sort'))
  mainElement.appendChild(header)

  fileInfos.forEach(({ Type, Name, Size, MTime }) => {
    const item = createElement('div', { class: 'DirView-fileInfo' })

    const detailsItem = createElement('div', { class: 'DirView-fileDetails' })
    // 类型图标
    const iconItem = createElement('div', { class: 'DirView-fileIcon' })
    iconItem.appendChild(createIcon(Type === 1 ? 'folder' : 'file'))
    detailsItem.appendChild(iconItem)
    // 文件名称
    const nameItem = createElement('div', { class: 'DirView-fileName' })
    nameItem.innerHTML = Name
    nameItem.onclick = () => { window.location.href = '/view' + path + (path.endsWith('/') ? '' : '/') + encodeURIComponent(Name) }
    detailsItem.appendChild(nameItem)
    // 折叠菜单
    const menuItem = createElement('div', { class: 'DirView-fileMenu' })
    menuItem.appendChild(createIcon('menu'))
    detailsItem.appendChild(menuItem)
    item.appendChild(detailsItem)

    // 文件大小
    const sizeItem = createElement('div', { class: 'DirView-fileSize' })
    sizeItem.innerHTML = calcFromBytes(Size)
    item.appendChild(sizeItem)

    // 修改时间
    const mtimeItem = createElement('div', {
      class: 'DirView-fileMTime',
      title: new Date(MTime * 1000).toLocaleString()
    })
    mtimeItem.innerHTML = calcRelativeTime(MTime * 1000)
    item.appendChild(mtimeItem)

    mainElement.appendChild(item)
  })

  if (fileInfos.length === 0) {
    const item = createElement('div', { class: 'DirView-fileInfo' })
    item.innerHTML = 'Sorry, this is an empty folder.'
    mainElement.appendChild(item)
  }

  return mainElement
}

export { createDirView }
