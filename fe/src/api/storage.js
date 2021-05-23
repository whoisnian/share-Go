import { fetchDeleteHead, fetchGetJSONWithStatus, fetchPostHead } from 'utils/request'

/** @enum { number } */
const FileType = {
  typeRegular: 0,
  typeDirectory: 1
}

/**
 * @typedef {{
 *   Type:  FileType,
 *   Name:  string,
 *   Size:  number,
 *   MTime: number
 * }} FileInfo
 */

/** @param { string } path */
const requestFileInfo = async (path) => {
  /** @type {{ ok: boolean, status: number, content: FileInfo | null }} */
  const result = await fetchGetJSONWithStatus(`/api/file${path}`)
  return result
}

/** @param { string } path */
const requestListDir = async (path) => {
  /** @type {{ ok: boolean, status: number, content: { FileInfos: FileInfo[] } | null }} */
  const result = await fetchGetJSONWithStatus(`/api/dir${path}`)
  return result
}

/** @param { string } path */
/** @param { FileList } files */
/** @param { Function } updateProgress */
const requestCreateFiles = async (path, files, updateProgress) => {
  for (let i = 0; i < files.length; i++) {
    const formData = new FormData()
    formData.append('fileList', files[i])

    await new Promise((resolve, reject) => {
      const xhr = new XMLHttpRequest()
      xhr.withCredentials = true
      xhr.onload = () => {
        if (200 <= xhr.status && xhr.status < 300) resolve(xhr.response)
        else reject(xhr.status)
      }
      xhr.onerror = () => reject(xhr.status)
      xhr.upload.onprogress = (event) => updateProgress(i, files.length, event.loaded, event.total)
      xhr.open('POST', `/api/upload${path}`)
      xhr.send(formData)
    })

    // await window.fetch(`/api/upload${path}`, {
    //   credentials: 'same-origin',
    //   method: 'POST',
    //   body: formData
    // })
  }
}

/** @param { string } path */
/** @param { string[] } urlList */
const requestDownloadFiles = async (path, urlList) => {
  const formData = new FormData()
  for (let i = 0; i < urlList.length; i++) {
    formData.append('urlList', urlList[i])
  }
  await window.fetch(`/api/upload${path}`, {
    credentials: 'same-origin',
    method: 'POST',
    body: formData
  })
}

/** @param { string } path */
const requestCreateDir = async (path) => {
  await fetchPostHead(`/api/dir${path}`)
}

/** @param { string } path */
const requestDeleteRecursively = async (path) => {
  await fetchDeleteHead(`/api/dir${path}`)
}

export {
  FileType,
  requestFileInfo,
  requestListDir,
  requestCreateFiles,
  requestDownloadFiles,
  requestCreateDir,
  requestDeleteRecursively
}
