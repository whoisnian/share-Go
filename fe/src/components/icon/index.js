import { createElement, createElementNS } from 'utils/element'
import { createMimeIcon } from './mime'

/**
 * @typedef {(
 *   'delete'|
 *   'download'|
 *   'edit'|
 *   'file-new'|
 *   'file'|
 *   'folder-new'|
 *   'folder-parent'|
 *   'folder'|
 *   'home'|
 *   'menu'|
 *   'paste'|
 *   'refresh'|
 *   'sort'|
 *   'sort-amount-asc'|
 *   'sort-amount-dsc'|
 *   'sort-name-asc'|
 *   'sort-name-dsc'|
 *   'sort-time-asc'|
 *   'sort-time-dsc'|
 *   'tab-new'
 * )} IconName
 */

/** @param { IconName } name */
const createIcon = (name, options = {}) => {
  const iconElement = createElement('div', options)
  const svgElement = createElementNS('http://www.w3.org/2000/svg', 'svg', {
    viewBox: '0 0 24 24',
    width: '24',
    height: '24'
  })

  if (name === 'delete') svgElement.innerHTML = '<path d="M9 4v2h1V5h4v1h1V4H9M5 7v1h14V7H5m2 2v11h10V9h-1v10H8V9H7"/>'
  else if (name === 'download') svgElement.innerHTML = '<path d="M9 4v6h1V5h4v5h1V4H9zm-2.207 7L6 11.816 12 18l6-6.184-.793-.816L12 16.367 6.793 11zM5 18v2h14v-2h-1v1H6v-1H5z"/>'
  else if (name === 'edit') svgElement.innerHTML = '<path d="M15.996 4L4.004 15.992H4V20h4.008v-.004L20 8.004l-.002-.002L20 8l-4-4-.002.002L15.996 4m-1.998 3.412l2.59 2.59L9 17.59V16H7V14.41l6.998-6.998M6 15.41V17h2V18.59l-.41.41H6l-1-1v-1.59l1-1"/>'
  else if (name === 'file-new') svgElement.innerHTML = '<path d="M5 4v16h9v-1H6V5h8v4h4v6h1V7.992L15.008 4 15 4.01V4H5zm11 11v2h-2v1h2v2h1v-2h2v-1h-2v-2h-1z"/>'
  else if (name === 'file') svgElement.innerHTML = '<path d="M5 20V4h10v.01l.008-.01L19 7.992V20H5zm1-1h12V9h-4V5H6v14z"/>'
  else if (name === 'folder-new') svgElement.innerHTML = '<path d="M4 4v16h10v-1H5v-8h3v-.01l.008.01 2-2H19v5h1V6h-6.992l-2-2-.008.01V4H4m13 11v2h-2v1h2v2h1v-2h2v-1h-2v-2h-1"/>'
  else if (name === 'folder-parent') svgElement.innerHTML = '<path d="M4 4v16h9v-1H5v-8h3v-.008l.008.008 2-2H19v4h1V6h-6.992l-2-2-.008.008V4H4zm12.5 8.793l-.707.707L13 16.293l.707.707L16 14.707V20h1v-5.293L19.293 17l.707-.707-2.793-2.793-.707-.707z"/>'
  else if (name === 'folder') svgElement.innerHTML = '<path d="M4 4v16h16V6h-6.992l-2-2-.008.008V4H4zm1 1h5.586l1.004 1.004L7.57 10H5V5zm5.008 4H19v10H5v-8h3v-.008l.008.008 2-2z"/>'
  else if (name === 'home') svgElement.innerHTML = '<path d="M12 4l-.707.707L4 12l.707.707.293-.293V20h14v-7.586l.293.293L20 12l-3-3V6h-3l-1.293-1.293L12 4m0 1.414l6 6V19h-4v-5h-4v5H6v-7.586l6-6"/>'
  else if (name === 'menu') svgElement.innerHTML = '<path d="M18 11h-2v2h2v-2m-5 0h-2v2h2v-2m-5 0H6v2h2v-2"/>'
  else if (name === 'paste') svgElement.innerHTML = '<path d="M8 4v2H5v14h14V6h-3V4H8zM6 7h1v2h10V7h1v12H6V7zm2 3v1h8v-1H8zm0 3v1h6v-1H8zm0 3v1h3v-1H8z"/>'
  else if (name === 'refresh') svgElement.innerHTML = '<path d="M20 12c0 1.442-.383 2.79-1.045 3.955l-.738-.738.002-.002-2.778-2.78.707-.706 2.481 2.482A6.946 6.946 0 0 0 19 12c0-3.878-3.122-7-7-7a6.985 6.985 0 0 0-3.217.783l-.738-.738A7.981 7.981 0 0 1 12 4c4.432 0 8 3.568 8 8zm-4.045 6.955A7.981 7.981 0 0 1 12 20c-4.432 0-8-3.568-8-8 0-1.442.383-2.79 1.045-3.955l.684.684.002-.002 2.828 2.828-.707.707-2.479-2.479A6.943 6.943 0 0 0 5 12c0 3.878 3.122 7 7 7a6.985 6.985 0 0 0 3.217-.783z"/>'
  else if (name === 'sort') svgElement.innerHTML = '<path d="M16 11v7.086l-2.293-2.293L13 16.5l3.5 3.5 3.5-3.5-.707-.707L17 18.086V11zM4 4v1h7V4zM4 7v1h5V7zM4 10v1h7v-1zM4 20v-1h7v1zM4 17v-1h7v1zM4 14v-1h7v1z"/><path d="M13 14V6.914l-2.293 2.293L10 8.5 13.5 5 17 8.5l-.707.707L14 6.914V14z"/>'
  else if (name === 'sort-amount-asc') svgElement.innerHTML = '<path d="M16 13V5.914l-2.293 2.293L13 7.5 16.5 4 20 7.5l-.707.707L17 5.914V13zM4 20v-1h14v1zM4 17v-1h12v1zM4 14v-1h10v1zM4 4v1h4V4zM4 7v1h6V7zM4 10v1h8v-1z"/>'
  else if (name === 'sort-amount-dsc') svgElement.innerHTML = '<path d="M16 11v7.086l-2.293-2.293L13 16.5l3.5 3.5 3.5-3.5-.707-.707L17 18.086V11zM4 4v1h14V4zM4 7v1h12V7zM4 10v1h10v-1zM4 20v-1h4v1zM4 17v-1h6v1zM4 14v-1h8v1z"/>'
  else if (name === 'sort-name-asc') svgElement.innerHTML = '<path d="M8.096 4c-.784.018-1.385.155-1.823.355l.247.9c.209-.072 1.002-.325 1.466-.325h.063c.646.027.693.262.693.818v.219l-.338.054C6.838 6.25 6 6.441 6 7.643 6 8.48 6.61 9 7.494 9c.61 0 1.01-.255 1.42-.602l.127.51h.967V5.54C10.008 4.51 9.296 4 8.203 4h-.107zm.646 2.678v1.119c-.337.228-.592.392-1.03.392-.39 0-.583-.245-.583-.591 0-.656.574-.765 1.275-.866l.338-.054zM14.5 9.793l-.207.207L11 13.293l.707.707L14 11.707V20h1v-8.293L17.293 14l.707-.707L14.707 10l-.207-.207zm-8.389.406v.815h2.073L6 14.27V15h3.828v-.852H7.582l2.174-3.29v-.659H6.11z"/>'
  else if (name === 'sort-name-dsc') svgElement.innerHTML = '<path d="M6.111 4v.814h2.073L6 8.068v.733h3.828v-.852H7.582l2.174-3.293V4H6.11zM14 9.992v8.293l-2.293-2.293L11 16.7l3.293 3.293.207.207.207-.207L18 16.7l-.707-.707L15 18.285V9.992h-1zM8.096 10c-.784.018-1.385.155-1.823.355l.247.9c.209-.072 1.002-.325 1.466-.325h.063c.646.027.693.262.693.818v.219l-.338.054C6.838 12.25 6 12.441 6 13.643 6 14.48 6.61 15 7.494 15c.61 0 1.01-.255 1.42-.602l.127.51h.967V11.54c0-1.03-.712-1.539-1.805-1.539h-.107zm.646 2.678v1.119c-.337.228-.592.392-1.03.392-.39 0-.583-.245-.583-.591 0-.656.574-.765 1.275-.866l.338-.054z"/>'
  else if (name === 'sort-time-asc') svgElement.innerHTML = '<path d="M20 12.002a.727.727 0 0 1-1.455 0 6.547 6.547 0 1 0-6.547 6.546.727.727 0 0 1 0 1.455A8.004 8.004 0 0 1 8.884 4.63 8 8 0 0 1 20 12.002zm-3.123.94a.727.727 0 0 0-1.029 0l-2.182 2.182a.728.728 0 0 0 1.029 1.03l.94-.941v4.063a.727.727 0 0 0 1.455 0v-4.063l.94.94a.728.728 0 0 0 1.03-1.029zm-4.879-6.055a.727.727 0 0 0-.727.727v4.388H8.36a.727.727 0 0 0 0 1.454h3.637a.727.727 0 0 0 .728-.727V7.614a.727.727 0 0 0-.728-.727z"/>'
  else if (name === 'sort-time-dsc') svgElement.innerHTML = '<path d="M20 12.001a.727.727 0 0 1-1.455 0 6.547 6.547 0 1 0-6.547 6.547.727.727 0 0 1 0 1.455A8.004 8.004 0 0 1 8.884 4.63 8 8 0 0 1 20 12zm-4.152 7.79a.727.727 0 0 0 1.03 0l2.182-2.183a.728.728 0 0 0-1.03-1.029l-.94.94v-4.063a.727.727 0 0 0-1.455 0v4.064l-.94-.94a.728.728 0 0 0-1.029 1.028zm-3.85-12.904a.727.727 0 0 0-.727.727v4.387H8.36a.727.727 0 0 0 0 1.455h3.637a.727.727 0 0 0 .728-.727V7.614a.727.727 0 0 0-.728-.727z"/>'
  else if (name === 'tab-new') svgElement.innerHTML = '<path d="M4 4v16h10v-1H5V5h4v3h10v6h1V7h-3V4H4zm13 11v2h-2v1h2v2h1v-2h2v-1h-2v-2h-1z"/>'
  else svgElement.innerHTML = '<path d="M11.969 4C10.218 4 8.562 4.422 7 5.264l.777 1.787c.781-.384 1.483-.658 2.108-.824a7.566 7.566 0 0 1 1.935-.25c.987 0 1.744.22 2.27.662.526.442.789 1.075.789 1.9 0 .442-.057.83-.172 1.164a3.198 3.198 0 0 1-.592 1c-.28.334-.867.883-1.763 1.65C11.252 13.358 10.1 15.14 10 17h2l-.006-.031c0-.759.132-1.368.395-1.826.271-.467.813-1.055 1.627-1.764.994-.842 1.657-1.475 1.986-1.9.337-.426.588-.875.752-1.35.164-.475.246-1.022.246-1.639 0-1.417-.448-2.519-1.344-3.303C14.76 4.396 13.531 4 11.97 4zM10 18v2h2v-2h-2z"/>'

  iconElement.appendChild(svgElement)
  return iconElement
}

export { createIcon, createMimeIcon }
