import { FileType, requestListDir, requestDeleteRecursively, requestCreateDir } from 'api/storage'
import { createIcon } from 'components/icon'
import { createContextMenu } from 'components/contextMenu'
import { createInputDialog } from 'components/inputDialog'
import { createUploadDialog } from 'components/uploadDialog'
import { createElement, downloadFile, copyText } from 'utils/element'
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
  pasteIcon.onclick = () => copyText(window.location.href)
  const folderNewIcon = createIcon('folder-new', { class: 'DirView-iconButton', title: 'Create new folder' })
  folderNewIcon.onclick = () => {
    const inputDialog = createInputDialog("Folder Name:", "new folder", (dirName) => {
      requestCreateDir(joinPath('/', oriPath, encodeURIComponent(dirName))).then(() => reloadPage())
    })
    const removeSelf = (event) => {
      if (event.target === inputDialog) {
        inputDialog.remove()
        document.removeEventListener('click', removeSelf)
      }
    }
    document.addEventListener('click', removeSelf)
    header.appendChild(inputDialog)
    inputDialog.focus()
  }
  const fileNewIcon = createIcon('file-new', { class: 'DirView-iconButton', title: 'Create new file' })
  fileNewIcon.onclick = () => {
    const { uploadDialog } = createUploadDialog(joinPath('/', oriPath))
    const removeSelf = (event) => {
      if (event.target === uploadDialog) {
        uploadDialog.remove()
        document.removeEventListener('click', removeSelf)
      }
    }
    document.addEventListener('click', removeSelf)
    header.appendChild(uploadDialog)
  }
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
    e.preventDefault()
    fileItem.showFileMenu(e)
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

  // Drag-and-Drop File Uploader
  main.ondragenter = main.ondragover = (e) => e.preventDefault()
  main.ondrop = (e) => {
    const { uploadDialog, uploadFiles } = createUploadDialog(joinPath('/', oriPath))
    const removeSelf = (event) => {
      if (event.target === uploadDialog) {
        uploadDialog.remove()
        document.removeEventListener('click', removeSelf)
      }
    }
    document.addEventListener('click', removeSelf)
    main.appendChild(uploadDialog)
    uploadFiles(e.dataTransfer.files)
    e.preventDefault()
  }

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
    listener: ({ data }) => copyText(window.location.origin + joinPath('/api/download', oriPath, encodeURIComponent(data.Name)))
  }, {
    icon: 'edit',
    name: '重命名',
    listener: () => alert('todo')
  }, {
    icon: 'download',
    name: '下载',
    listener: ({ data }) => downloadFile(joinPath('/api/download', oriPath, encodeURIComponent(data.Name)), data.FileType === FileType.typeDirectory ? data.Name + '.zip' : data.Name)
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
    fileItem.showFileMenu = (event) => showFileMenu(event, info)
    fileItem.oncontextmenu = (event) => {
      event.preventDefault()
      fileItem.showFileMenu(event)
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
