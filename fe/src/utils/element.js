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

export { getRootElement }
