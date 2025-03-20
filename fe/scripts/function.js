import { copyFileSync, existsSync, lstatSync, mkdirSync, readdirSync } from 'fs'
import { resolve, join } from 'path'

const PATH_ROOT = resolve(import.meta.dirname, '..')
const PATH_OUTPUT = resolve(import.meta.dirname, '../dist')
const fromRoot = (...args) => resolve(PATH_ROOT, ...args)
const fromOutput = (...args) => resolve(PATH_OUTPUT, ...args)

/**
 * @param { string } src
 * @param { string } dest
 */
const copyRecursivelySync = (src, dest) => {
  if (!src || !dest || !existsSync(src)) return

  const fileStat = lstatSync(src)
  if (fileStat.isFile()) {
    copyFileSync(src, dest)
  } else if (fileStat.isDirectory()) {
    if (!existsSync(dest)) mkdirSync(dest, { recursive: true })
    const files = readdirSync(src)
    files.forEach((file) => {
      copyRecursivelySync(join(src, file), join(dest, file))
    })
  }
}

export {
  fromRoot,
  fromOutput,
  copyRecursivelySync
}
