import { fetchGetJSON, fetchGetJSONWithStatus } from 'utils/request'

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
  /** @type {{ ok: boolean, status: number, content: FileInfo }} */
  const result = await fetchGetJSONWithStatus(`/api/file${path}`)
  return result
}

/** @param { string } path */
const requestListDir = async (path) => {
  /** @type {{ FileInfos: FileInfo[] }} */
  const { FileInfos } = await fetchGetJSON(`/api/dir${path}`)
  return FileInfos
}

export {
  FileType,
  requestFileInfo,
  requestListDir
}