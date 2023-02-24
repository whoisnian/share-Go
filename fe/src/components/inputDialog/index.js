import { createElement } from 'utils/element'
import './style.css'

/**
 * @param { HTMLElement } parent
 * @param { string } name
 * @param { string } preset
 * @param { function } listener
 */
const createInputDialog = (parent, name, preset, listener) => {
  const inputDialog = createElement('div', { class: 'InputDialog-modal' })
  const popup = createElement('div', { class: 'InputDialog-popup' })
  const header = createElement('div', { class: 'InputDialog-header' })
  header.textContent = name
  const input = createElement('input', { class: 'InputDialog-input', type: 'text' })
  input.value = preset
  inputDialog.focus = () => {
    input.focus()
    input.select()
  }
  const button = createElement('div', { class: 'InputDialog-button' })
  button.textContent = 'OK'

  const keyCheck = (event) => {
    if (event.target === input && event.key === 'Enter') button.click()
  }
  const removeSelf = (event) => {
    if (event.target === inputDialog || event.target === button) {
      document.removeEventListener('click', removeSelf)
      input.removeEventListener('keypress', keyCheck)
      inputDialog.remove()
    }
  }
  document.addEventListener('click', removeSelf)
  input.addEventListener('keypress', keyCheck)
  button.onclick = (event) => {
    listener(input.value)
    removeSelf(event)
  }

  popup.appendChild(header)
  popup.appendChild(input)
  popup.appendChild(button)
  inputDialog.appendChild(popup)
  parent.appendChild(inputDialog)

  inputDialog.focus()
  return inputDialog
}

export { createInputDialog }
