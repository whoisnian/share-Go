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

const chooseFile = async (multiple) => {
  const input = createElement('input', { type: 'file', multiple })
  return new Promise((resolve) => {
    input.oncancel = () => {
      input.remove()
      resolve([])
    }
    input.onchange = () => {
      const files = input.files
      input.remove()
      resolve(files)
    }
    input.click()
  })
}

const copyText = (text) => {
  if (window.navigator && window.isSecureContext) {
    window.navigator.clipboard.writeText(text)
  } else {
    const input = createElement('input', { type: 'text', style: 'position:fixed;', value: text })
    document.body.appendChild(input)
    input.select()
    document.execCommand('copy')
    input.remove()
  }
}

export {
  getRootElement,
  createElement,
  createElementNS,
  downloadFile,
  chooseFile,
  copyText
}
