import { getRootElement } from 'utils/element'

const init404Page = async () => {
  getRootElement()
    .removeAllChildren()
    .appendChild(document.createTextNode('Page Not Found'))
}

export { init404Page }
