import { createElement, chooseFile } from 'utils/element'
import { createInfoDialog } from 'components/infoDialog'
import { requestCreateFile, requestUploadFiles, requestDownloadFiles } from 'api/storage'
import { joinPath, reloadPage } from 'utils/function'
import './style.css'

/**
 * @param { HTMLElement } parent
 * @param { string } base
 */
const createUploadDialog = (parent, base) => {
  const uploadDialog = createElement('div', { class: 'UploadDialog-modal' })
  const popup = createElement('div', { class: 'UploadDialog-popup' })
  const header = createElement('div', { class: 'UploadDialog-header' })
  header.textContent = 'Upload From:'

  const tabList = []
  const tab = createElement('div', { class: 'UploadDialog-tab' })
  tab.changeTo = (child) => {
    if (tab.current !== child) {
      if (tab.current) {
        tab.current.className = tab.current.className.replace(' active ', '')
        tab.current.tabContent.style.display = 'none'
      }
      child.className += ' active '
      child.tabContent.style.display = 'flex'
      tab.current = child
    }
  }
  tab.appendChildWithTabContent = (child) => {
    tabList.push(child)
    tab.appendChild(child)
    child.onclick = () => tab.changeTo(child)
  }

  // 从本地选择文件上传
  const fromFiles = createElement('div', { class: 'UploadDialog-tabItem' })
  fromFiles.textContent = 'Files'
  fromFiles.tabContent = createElement('div', { class: 'UploadDialog-tabContent' })
  const filesButton = createElement('label', { class: 'UploadDialog-button' })
  filesButton.textContent = 'Browse'
  const filesContent = createElement('input', { class: 'UploadDialog-input', type: 'text', placeholder: 'No files choosed', readonly: true })
  const updateProgress = (fileIndex, fileTotal, dataLoaded, dataTotal) => {
    filesButton.textContent = dataTotal ? `${Math.round(100 * dataLoaded / dataTotal)} %` : 'Uploading'
    filesContent.value = `Uploading ${fileIndex + 1} of ${fileTotal} files.`
  }
  const uploadFilesThenReload = (fileList) => {
    if (fileList?.length <= 0) {
      return false
    }

    filesButton.textContent = 'Ready'
    if (fileList.length === 1) {
      filesContent.value = fileList[0].name
    } else if (fileList.length > 1) {
      filesContent.value = `Choosed ${fileList.length} files.`
    }
    requestUploadFiles(base, fileList, updateProgress).then(reloadPage).catch(console.error)
    return true
  }

  // 输入 url 从远端下载
  const fromUrl = createElement('div', { class: 'UploadDialog-tabItem' })
  fromUrl.textContent = 'Url'
  fromUrl.tabContent = createElement('div', { class: 'UploadDialog-tabContent' })
  const urlInput = createElement('input', { class: 'UploadDialog-input', type: 'text', placeholder: 'http(s)://example.com' })
  const urlButton = createElement('div', { class: 'UploadDialog-button' })
  urlButton.textContent = 'OK'
  const downloadFilesThenReload = () => {
    if (urlInput.value?.length <= 0) {
      createInfoDialog(popup, 'Error', 'missing url input')
      return false
    } else if (!urlInput.value.match(/^https?:\/\//)) {
      createInfoDialog(popup, 'Error', 'invalid url scheme')
      return false
    }

    const urlList = urlInput.value.split(/[,，]/)
    requestDownloadFiles(base, urlList).then(reloadPage).catch(console.error)
    return true
  }

  // 输入文本内容上传
  const fromText = createElement('div', { class: 'UploadDialog-tabItem' })
  fromText.textContent = 'Text'
  fromText.tabContent = createElement('div', { class: 'UploadDialog-tabContent', style: 'flex-direction:column;' })
  const textTitleButtonGroup = createElement('div', { class: 'UploadDialog-group' })
  const textTitle = createElement('input', { class: 'UploadDialog-input', type: 'text', value: 'untitled.txt' })
  const textButton = createElement('div', { class: 'UploadDialog-button' })
  textButton.textContent = 'Upload'
  const textContent = createElement('textarea', { class: 'UploadDialog-input', style: 'resize:vertical;', rows: 5, placeholder: 'Input text content here...' })
  const uploadTextThenReload = () => {
    if (textTitle.value?.length <= 0) {
      createInfoDialog(popup, 'Error', 'missing text title')
      return false
    }

    const filePath = joinPath(base, encodeURIComponent(textTitle.value))
    requestCreateFile(filePath, textContent.value).then(reloadPage).catch(console.error)
    return true
  }

  const keyCheck = (event) => {
    event.stopPropagation()
    if (event.target === urlInput && event.key === 'Enter') {
      event.target.blur()
      urlButton.click()
    } else if ((event.target === textTitle || event.target === textContent) && event.key === 'Enter') {
      event.target.blur()
      textButton.click()
    }
  }
  const removeSelf = (event) => {
    if (event.target === uploadDialog) {
      document.removeEventListener('click', removeSelf)
      urlInput.removeEventListener('keypress', keyCheck)
      textTitle.removeEventListener('keypress', keyCheck)
      textContent.removeEventListener('keypress', keyCheck)
      uploadDialog.remove()
    }
  }
  document.addEventListener('click', removeSelf)
  urlInput.addEventListener('keypress', keyCheck)
  textTitle.addEventListener('keypress', keyCheck)
  textContent.addEventListener('keypress', keyCheck)
  filesButton.onclick = async () => {
    const files = await chooseFile(true)
    if (uploadFilesThenReload(files)) document.removeEventListener('click', removeSelf)
  }
  urlButton.onclick = () => {
    if (downloadFilesThenReload()) document.removeEventListener('click', removeSelf)
  }
  textButton.onclick = () => {
    if (uploadTextThenReload()) document.removeEventListener('click', removeSelf)
  }

  fromFiles.tabContent.appendChild(filesButton)
  fromFiles.tabContent.appendChild(filesContent)
  fromUrl.tabContent.appendChild(urlInput)
  fromUrl.tabContent.appendChild(urlButton)
  textTitleButtonGroup.appendChild(textTitle)
  textTitleButtonGroup.appendChild(textButton)
  fromText.tabContent.appendChild(textTitleButtonGroup)
  fromText.tabContent.appendChild(textContent)
  tab.appendChildWithTabContent(fromFiles)
  tab.appendChildWithTabContent(fromUrl)
  tab.appendChildWithTabContent(fromText)
  tab.changeTo(fromFiles)

  popup.appendChild(header)
  popup.appendChild(tab)
  tabList.forEach(t => popup.appendChild(t.tabContent))
  uploadDialog.appendChild(popup)
  parent.appendChild(uploadDialog)

  return {
    uploadDialog,
    uploadFilesThenReload
  }
}

export { createUploadDialog }
