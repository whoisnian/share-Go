/** @returns { HTMLElement } */
const getRootElement = () => {
  const rootElement = document.getElementById('root')
  rootElement.removeAllChildren = () => {
    while (rootElement.firstChild) {
      rootElement.removeChild(rootElement.firstChild)
    }
    return rootElement
  }
  return rootElement
}

/** @returns { HTMLElement } */
const createElement = (tag, options = {}) => {
  const element = document.createElement(tag)
  Object.entries(options).forEach(([k, v]) => {
    element.setAttribute(k, v)
  })
  return element
}

/** @returns { HTMLElement } */
const createElementNS = (namespace, tag, options = {}) => {
  const element = document.createElementNS(namespace, tag)
  Object.entries(options).forEach(([k, v]) => {
    element.setAttribute(k, v)
  })
  return element
}

const downloadFile = (url, filename) => {
  const link = createElement('a', {
    href: url,
    download: filename
  })
  link.click()
  link.remove()
}

const chooseFile = (listener, multiple) => {
  const input = createElement('input', { type: 'file', multiple })
  input.onchange = () => {
    listener(input.files)
    input.remove()
  }
  input.click()
}

export {
  getRootElement,
  createElement,
  createElementNS,
  downloadFile,
  chooseFile
}
