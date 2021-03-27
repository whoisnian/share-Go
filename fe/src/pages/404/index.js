import { getRootElement, removeAllChildren } from 'utils/element'

const init404Page = async () => {
  const rootElement = getRootElement()
  removeAllChildren(rootElement)
  rootElement.appendChild(document.createTextNode('404 Not Found'))
}

export { init404Page }
