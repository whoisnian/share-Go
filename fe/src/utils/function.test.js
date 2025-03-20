import test from 'node:test'
import { strictEqual } from 'assert'
import { joinPath } from './function.js'

test('test joinPath', async (t) => {
  const tests = [
    [['a', 'b', 'c'], 'a/b/c'],
    [['/a', 'b/c'], '/a/b/c'],
    [['/./a', 'b/c/.'], '/a/b/c'],
    [['//////', './.'], '/'],
    [['a/b', '../c'], 'a/c'],
    [['a/b', '../../../c'], '../c'],
    [['/a/b', '../../../..'], '/']
  ]

  for (let i = 0; i < tests.length; i++) {
    const [input, output] = tests[i]
    const actual = joinPath(...input)
    await t.test(`subtest ${i}`, () => { strictEqual(actual, output) })
  }
})
