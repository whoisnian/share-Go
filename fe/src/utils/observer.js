/** @type { Map<string, any> } */
const objMap = new Map()

/** @type { Map<string, Function[]> } */
const cbMap = new Map()

/**
 * @typedef {(
 *   'pathname'|
 *   'page'|
 *   'viewPage'|
 *   'oriPath'
 * )} ObserverKey
 */

const observer = {
  /** @param { ObserverKey } k */
  set: (k, v) => {
    const prev = objMap.get(k)
    if (prev !== v) {
      objMap.set(k, v)
      cbMap.get(k).forEach(cb => cb(prev, v))
    }
  },
  /** @param { ObserverKey } k */
  get: (k) => {
    return objMap.get(k)
  },
  /** @param { ObserverKey } k */
  delete: (k) => {
    const prev = objMap.get(k)
    if (prev !== undefined) {
      objMap.delete(k)
      cbMap.get(k).forEach(cb => cb(prev, undefined))
    }
  },
  /**
   * @param { ObserverKey } k
   * @param { Function } cb
   */
  onchange: (k, cb) => {
    if (cbMap.has(k)) cbMap.get(k).push(cb)
    else cbMap.set(k, [cb])
  }
}

export { observer }
