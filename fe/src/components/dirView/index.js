import { requestListDir } from 'api/storage'
import { createIcon } from 'components/icon'
import { createContextMenu } from 'components/contextMenu'
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
    return b.Type - a.Type
  })

  const mainElement = createElement('div', { class: 'DirView-main' })
  const header = createElement('div', { class: 'DirView-header' })
  const currentPath = createElement('div', { class: 'DirView-currentPath' })
  currentPath.textContent = decodeURIComponent(path)

  const parentIcon = createIcon('folder-parent', {
    class: 'DirView-iconButton',
    title: 'Go to parent folder'
  })
  header.appendChild(parentIcon)

  const refreshIcon = createIcon('refresh', {
    class: 'DirView-iconButton',
    title: 'Refresh'
  })
  header.appendChild(refreshIcon)

  const homeIcon = createIcon('home', {
    class: 'DirView-iconButton',
    title: 'Go to home'
  })
  header.appendChild(homeIcon)

  header.appendChild(currentPath)

  const pasteIcon = createIcon('paste', {
    class: 'DirView-iconButton',
    title: 'Copy current path'
  })
  header.appendChild(pasteIcon)

  const folderNewIcon = createIcon('folder-new', {
    class: 'DirView-iconButton',
    title: 'Create new folder'
  })
  header.appendChild(folderNewIcon)

  const fileNewIcon = createIcon('file-new', {
    class: 'DirView-iconButton',
    title: 'Create new file'
  })
  header.appendChild(fileNewIcon)

  const sortIcon = createIcon('sort', {
    class: 'DirView-iconButton',
    title: 'Sort by'
  })
  header.appendChild(sortIcon)

  mainElement.appendChild(header)

  const {
    contextMenu: fileMenu,
    show: showFileMenu
  } = createContextMenu()
  mainElement.appendChild(fileMenu)

  fileInfos.forEach(({ Type, Name, Size, MTime }) => {
    const item = createElement('div', { class: 'DirView-fileInfo' })
    item.oncontextmenu = (e) => {
      e.preventDefault()
      showFileMenu(e)
    }

    const detailsItem = createElement('div', { class: 'DirView-fileDetails' })
    // 类型图标
    const iconItem = createIcon(Type === 1 ? 'folder' : 'file', { class: 'DirView-fileIcon' })
    detailsItem.appendChild(iconItem)
    // 文件名称
    const nameItem = createElement('span', { class: 'DirView-fileName' })
    const nameLink = createElement('a', {
      class: 'DirView-nameLink',
      title: Name,
      href: '/view' + path + (path.endsWith('/') ? '' : '/') + encodeURIComponent(Name)
    })
    nameLink.textContent = Name
    nameItem.appendChild(nameLink)
    detailsItem.appendChild(nameItem)
    // 折叠菜单
    const menuItem = createIcon('menu', { class: 'DirView-iconButton DirView-fileMenu' })
    menuItem.onclick = showFileMenu
    detailsItem.appendChild(menuItem)
    item.appendChild(detailsItem)

    // 文件大小
    const sizeItem = createElement('div', { class: 'DirView-fileSize' })
    sizeItem.textContent = calcFromBytes(Size)
    item.appendChild(sizeItem)

    // 修改时间
    const mtimeItem = createElement('div', {
      class: 'DirView-fileMTime',
      title: new Date(MTime * 1000).toLocaleString()
    })
    mtimeItem.textContent = calcRelativeTime(MTime * 1000)
    item.appendChild(mtimeItem)

    mainElement.appendChild(item)
  })

  if (fileInfos.length === 0) {
    const item = createElement('div', { class: 'DirView-fileInfo' })
    item.textContent = 'Sorry, this is an empty folder.'
    mainElement.appendChild(item)
  }

  return mainElement
}

export { createDirView }
