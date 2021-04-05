import { createElement } from 'utils/element'
import './style.css'

const createContextMenu = () => {
  const contextMenu = createElement('div', { class: 'ContextMenu-menu' })

  document.addEventListener('click', () => { contextMenu.style.display = 'none' })

  const show = (event) => {
    event.cancelBubble = true
    contextMenu.style.display = 'flex'
    contextMenu.style.top = event.clientY + 'px'
    contextMenu.style.left = event.clientX + 'px'
    contextMenu.textContent = 'This is a contextMenu.'
    console.log(event)
  }

  return {
    contextMenu,
    show
  }
}

export { createContextMenu }
