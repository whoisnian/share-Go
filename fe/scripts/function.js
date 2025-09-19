import { resolve } from 'path'

const PATH_ROOT = resolve(import.meta.dirname, '..')
const PATH_OUTPUT = resolve(import.meta.dirname, '../dist')
const fromRoot = (...args) => resolve(PATH_ROOT, ...args)
const fromOutput = (...args) => resolve(PATH_OUTPUT, ...args)

export {
  fromRoot,
  fromOutput
}
