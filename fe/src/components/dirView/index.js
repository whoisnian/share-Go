import { FileType, requestListDir, requestDeleteRecursively } from 'api/storage'
import { createIcon } from 'components/icon'
import { createContextMenu } from 'components/contextMenu'
import { createElement, downloadFile } from 'utils/element'
import { calcFromBytes, calcRelativeTime, joinPath, openUrl, openUrlInNewTab, reloadPage } from 'utils/function'
import './style.css'

/** @param { string } oriPath */
const createHeader = (oriPath) => {
  const currentPath = joinPath('/', decodeURIComponent(oriPath))
  const header = createElement('div', { class: 'DirView-header' })

  const parentIcon = createIcon('folder-parent', { class: 'DirView-iconButton', title: 'Go to parent folder' })
  parentIcon.onclick = () => openUrl('/view' + joinPath('/', oriPath, '..'))
  const refreshIcon = createIcon('refresh', { class: 'DirView-iconButton', title: 'Refresh' })
  refreshIcon.onclick = () => reloadPage()
  const homeIcon = createIcon('home', { class: 'DirView-iconButton', title: 'Go to home' })
  homeIcon.onclick = () => openUrl('/view/')
  const pathSpan = createElement('span', { class: 'DirView-pathSpan', title: currentPath })
  pathSpan.textContent = currentPath
  const pasteIcon = createIcon('paste', { class: 'DirView-iconButton', title: 'Copy current url' })
  pasteIcon.onclick = () => navigator.clipboard.writeText(window.location.href)
  const folderNewIcon = createIcon('folder-new', { class: 'DirView-iconButton', title: 'Create new folder' })
  const fileNewIcon = createIcon('file-new', { class: 'DirView-iconButton', title: 'Create new file' })
  const sortIcon = createIcon('sort', { class: 'DirView-iconButton', title: 'Sort by' })

  header.appendChild(parentIcon)
  header.appendChild(refreshIcon)
  header.appendChild(homeIcon)
  header.appendChild(pathSpan)
  header.appendChild(pasteIcon)
  header.appendChild(folderNewIcon)
  header.appendChild(fileNewIcon)
  header.appendChild(sortIcon)
  return header
}

/** @param { string } oriPath */
/** @param { import('api/storage').FileInfo } fileInfo */
const createFileItem = (oriPath, fileInfo) => {
  const fileItem = createElement('div', { id: fileInfo.Name, class: 'DirView-fileInfo' })

  // 文件详情（图标，名称，菜单）
  const detailsItem = createElement('div', { class: 'DirView-fileDetails' })
  const iconItem = createIcon(fileInfo.Type === 1 ? 'folder' : 'file', { class: 'DirView-fileIcon' })
  const nameItem = createElement('span', { class: 'DirView-fileName' })
  const nameLink = createElement('a', {
    class: 'DirView-nameLink',
    title: fileInfo.Name,
    href: joinPath('/view', oriPath, encodeURIComponent(fileInfo.Name))
  })
  nameLink.textContent = fileInfo.Name
  nameItem.appendChild(nameLink)
  const menuItem = createIcon('menu', { class: 'DirView-iconButton DirView-fileMenu' })
  menuItem.onclick = (e) => {
    e.cancelBubble = true
    fileItem.dispatchEvent(new window.MouseEvent('contextmenu', {
      clientX: e.clientX,
      clientY: e.clientY,
      button: 2,
      buttons: 2
    }))
  }
  detailsItem.appendChild(iconItem)
  detailsItem.appendChild(nameItem)
  detailsItem.appendChild(menuItem)
  // 文件大小
  const sizeItem = createElement('div', { class: 'DirView-fileSize' })
  sizeItem.textContent = calcFromBytes(fileInfo.Size)
  // 修改时间
  const mtimeItem = createElement('div', {
    class: 'DirView-fileMTime',
    title: new Date(fileInfo.MTime * 1000).toLocaleString()
  })
  mtimeItem.textContent = calcRelativeTime(fileInfo.MTime * 1000)

  fileItem.appendChild(detailsItem)
  fileItem.appendChild(sizeItem)
  fileItem.appendChild(mtimeItem)
  return fileItem
}

/** @param { string } oriPath */
const createDirView = async (oriPath) => {
  const { ok, status, content } = await requestListDir(oriPath)
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

  const main = createElement('div', { class: 'DirView-main' })

  const header = createHeader(oriPath)
  main.appendChild(header)

  const {
    contextMenu: fileMenu,
    show: showFileMenu
  } = createContextMenu([{
    icon: 'tab-new',
    name: '新建标签页打开',
    listener: ({ data }) => openUrlInNewTab(joinPath('/view', oriPath, encodeURIComponent(data.Name)))
  }, {
    icon: 'paste',
    name: '复制下载链接',
    listener: ({ data }) => navigator.clipboard.writeText(window.location.origin + joinPath('/download', oriPath, encodeURIComponent(data.Name)))
  }, {
    icon: 'edit',
    name: '重命名',
    listener: () => console.log('todo')
  }, {
    icon: 'download',
    name: '下载',
    listener: ({ data }) => downloadFile(joinPath('/download', oriPath, encodeURIComponent(data.Name)), data.FileType === FileType.typeDirectory ? data.Name + '.zip' : data.Name)
  }, {
    icon: 'delete',
    name: '删除',
    listener: ({ event, data }) => {
      requestDeleteRecursively(joinPath(oriPath, encodeURIComponent(data.Name)))
      let item = event.target
      while (item.parentElement && (item.className !== 'DirView-fileInfo' || item.id !== data.Name)) item = item.parentElement
      if (item.className === 'DirView-fileInfo' && item.id === data.Name) item.remove()
    }
  }])
  main.appendChild(fileMenu)

  fileInfos.forEach(info => {
    const fileItem = createFileItem(oriPath, info)
    fileItem.oncontextmenu = (event) => {
      event.preventDefault()
      showFileMenu(event, info)
    }
    main.appendChild(fileItem)
  })

  if (fileInfos.length === 0) {
    const item = createElement('div', { class: 'DirView-fileInfo' })
    item.textContent = 'Sorry, this is an empty folder.'
    main.appendChild(item)
  }

  return main
}

export { createDirView }
