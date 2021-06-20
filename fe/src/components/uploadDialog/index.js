import { createElement, chooseFile } from 'utils/element'
import { requestCreateFiles, requestDownloadFiles } from 'api/storage'
import { reloadPage } from 'utils/function'
import './style.css'

/**
 * @param { string } name
 * @param { string } preset
 * @param { function } listener
 */
const createUploadDialog = (base) => {
  const uploadDialog = createElement('div', { class: 'UploadDialog-modal' })
  const popup = createElement('div', { class: 'UploadDialog-popup' })
  const header = createElement('div', { class: 'UploadDialog-header' })
  header.textContent = 'Upload:'

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
  const fromLocal = createElement('div', { class: 'UploadDialog-tabItem' })
  fromLocal.textContent = 'From Local'
  fromLocal.tabContent = createElement('div', { class: 'UploadDialog-tabContent' })
  const uploadButton = createElement('label', { class: 'UploadDialog-button' })
  uploadButton.textContent = 'Browse'
  const uploadContent = createElement('input', { class: 'UploadDialog-input', type: 'text', readonly: true })
  const updateProgress = (fileIndex, fileTotal, dataLoaded, dataTotal) => {
    uploadButton.textContent = dataTotal ? `${Math.round(100 * dataLoaded / dataTotal)} %` : `Uploading`
    uploadContent.value = `Uploading ${fileIndex + 1} of ${fileTotal} files.`
  }
  // TODO: 修改 export 该方法的实现方式
  uploadDialog.uploadFiles = (fileList) => {
    uploadButton.textContent = 'Ready'
    if (fileList && fileList.length === 1) {
      uploadContent.value = fileList[0].name
    } else if (fileList && fileList.length > 1) {
      uploadContent.value = `Choosed ${fileList.length} files.`
    }
    requestCreateFiles(base, fileList, updateProgress).then(() => {
      uploadDialog.remove()
      reloadPage()
    }).catch((e) => {
      console.error(e)
      uploadDialog.remove()
    })
  }
  uploadButton.onclick = () => {
    chooseFile(uploadDialog.uploadFiles, true)
  }
  fromLocal.tabContent.appendChild(uploadButton)
  fromLocal.tabContent.appendChild(uploadContent)

  // 输入 url 从远端下载
  const fromUrl = createElement('div', { class: 'UploadDialog-tabItem' })
  fromUrl.textContent = 'From Url'
  fromUrl.tabContent = createElement('div', { class: 'UploadDialog-tabContent' })
  const input = createElement('input', { class: 'UploadDialog-input', type: 'text' })
  const button = createElement('div', { class: 'UploadDialog-button' })
  button.textContent = 'OK'
  button.onclick = () => {
    if (input.value?.length > 0) {
      const urlList = input.value.split(/[,，]/)
      requestDownloadFiles(base, urlList).then(() => {
        uploadDialog.remove()
        reloadPage()
      }).catch((e) => {
        console.error(e)
        uploadDialog.remove()
      })
    }
  }
  fromUrl.tabContent.appendChild(input)
  fromUrl.tabContent.appendChild(button)

  tab.appendChildWithTabContent(fromLocal)
  tab.appendChildWithTabContent(fromUrl)
  tab.changeTo(fromLocal)

  popup.appendChild(header)
  popup.appendChild(tab)
  tabList.forEach(t => popup.appendChild(t.tabContent))
  uploadDialog.appendChild(popup)
  return uploadDialog
}

export { createUploadDialog }
