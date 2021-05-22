const { copyFileSync, existsSync, lstatSync, mkdirSync, readdirSync } = require('fs')
const { resolve, join } = require('path')

const PATH_ROOT = resolve(__dirname, '..')
const PATH_OUTPUT = resolve(__dirname, '../dist')
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

module.exports = {
  fromRoot,
  fromOutput,
  copyRecursivelySync
}