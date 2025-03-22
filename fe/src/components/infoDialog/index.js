import { createElement } from 'utils/element'
import './style.css'

/**
 * @param { HTMLElement } parent
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

  const keyCheck = (event) => {
    if (event.key === 'Enter' && !event.ctrlKey && !event.altKey && !event.shiftKey) button.click()
  }
  const removeSelf = (event) => {
    if (event.target === infoDialog || event.target === button) {
      document.removeEventListener('click', removeSelf)
      document.removeEventListener('keypress', keyCheck)
      infoDialog.remove()
    }
  }
  document.addEventListener('click', removeSelf)
  document.addEventListener('keypress', keyCheck)
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
