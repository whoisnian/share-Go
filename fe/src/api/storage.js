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
const requestCreateFiles = async (path, files) => {
  const formData = new FormData()
  for (let i = 0; i < files.length; i++) {
    formData.append('fileList', files[i])
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
  requestCreateDir,
  requestDeleteRecursively
}
