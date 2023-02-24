import { createElement } from 'utils/element'
import './style.css'

/**
 * @param { string } title
 * @param { string } content
 */
const createInfoDialog = (parent, title, content) => {
  const infoDialog = createElement('div', { class: 'InfoDialog-modal' })
  const popup = createElement('div', { class: 'InfoDialog-popup' })
  const header = createElement('div', { class: 'InfoDialog-header' })
  header.textContent = title
  const body = createElement('div', { class: 'InfoDialog-body' })
  body.textContent = content
  const button = createElement('div', { class: 'InfoDialog-button' })
  button.textContent = 'OK'

  const isCloseEvent = (event) => {
    switch (event.type) {
      case 'click': return (event.target === infoDialog || event.target === button)
      case 'keypress': return event.key === 'Enter'
      default: return false
    }
  }
  const removeSelf = (event) => {
    if (isCloseEvent(event)) {
      infoDialog.remove()
      document.removeEventListener('click', removeSelf)
      document.removeEventListener('keypress', removeSelf)
    }
  }
  document.addEventListener('click', removeSelf)
  document.addEventListener('keypress', removeSelf)
  button.onclick = removeSelf

  popup.appendChild(header)
  popup.appendChild(body)
  popup.appendChild(button)
  infoDialog.appendChild(popup)
  parent.appendChild(infoDialog)

  infoDialog.focus()
  return infoDialog
}

export { createInfoDialog }
