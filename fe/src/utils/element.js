const getRootElement = () => {
  return document.getElementById('root')
}

const removeAllChildren = (element) => {
  while (element.firstChild) {
    element.removeChild(element.firstChild)
  }
}

export {
  getRootElement,
  removeAllChildren
}
