import { createIcon } from 'components/icon'
import { createElement } from 'utils/element'
import './style.css'

/**
 * @param {{
 *   icon: import('components/icon').IconName,
 *   name: string,
 *   listener: function
 * }[]} items
 */
const createContextMenu = (items) => {
  const contextMenu = createElement('div', { class: 'ContextMenu-menu' })
  items.forEach(({ icon, name, listener }) => {
    const menuItem = createElement('div', { class: 'ContextMenu-menuItem' })
    const itemIcon = createIcon(icon, { class: 'ContextMenu-itemIcon' })
    const itemName = createElement('span', { class: 'ContextMenu-itemName' })
    itemName.textContent = name

    menuItem.appendChild(itemIcon)
    menuItem.appendChild(itemName)
    menuItem.onclick = listener
    contextMenu.appendChild(menuItem)
  })

  document.addEventListener('click', () => { contextMenu.style.display = 'none' })

  const show = (event) => {
    event.cancelBubble = true
    contextMenu.style.display = 'flex'
    contextMenu.style.top = event.clientY + 'px'
    contextMenu.style.left = event.clientX + 'px'
    console.log(event)
  }

  return {
    contextMenu,
    show
  }
}

export { createContextMenu }
