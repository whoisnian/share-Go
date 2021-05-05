import { createFileView } from 'components/fileView'
import { createDirView } from 'components/dirView'
import { FileType, requestFileInfo } from 'api/storage'
import { observer } from 'utils/observer'

const initViewPage = async () => {
  const changeViewPage = async (_, oriPath) => {
    const { ok, status, content: fileInfo } = await requestFileInfo(oriPath)
    if (ok && fileInfo.Type === FileType.typeRegular) {
      observer.set('viewPage', createFileView)
    } else if (ok && fileInfo.Type === FileType.typeDirectory) {
      observer.set('viewPage', createDirView)
    } else {
      observer.set('viewPage', async () => document.createTextNode(ok
        ? 'Can not parse fileInfo'
        : status === 404
          ? 'File or dir not found'
          : 'Unexpected error'
      ))
    }
  }
  observer.onchange('oriPath', changeViewPage)

  let mainElement = document.createTextNode('loading1...')
  const render = async (_, viewPage) => {
    const newElement = await viewPage()
    mainElement.parentElement.replaceChild(newElement, mainElement)
    mainElement = newElement
  }
  observer.onchange('viewPage', render)

  const pathname = observer.get('pathname')
  const oriPath = pathname.slice(pathname.indexOf('/view') + 5)
  observer.set('oriPath', oriPath)

  return mainElement
}

export { initViewPage }
