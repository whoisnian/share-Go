import { createElement } from 'utils/element'
import './style.css'

/**
 * @param { string } name
 * @param { string } preset
 * @param { function } listener
 */
const createInputDialog = (name, preset, listener) => {
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
  button.onclick = () => {
    listener(input.value)
    inputDialog.remove()
  }

  popup.appendChild(header)
  popup.appendChild(input)
  popup.appendChild(button)
  inputDialog.appendChild(popup)
  return inputDialog
}

export { createInputDialog }
