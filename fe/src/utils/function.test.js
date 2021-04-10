import { strictEqual } from 'assert'
import { joinPath } from 'utils/function'

const testJoinPath = () => {
  console.log('testJoinPath start:')
  let cnt = 0
  let ecnt = 0

  const tests = [
    [['a', 'b', 'c'], 'a/b/c'],
    [['/a', 'b/c'], '/a/b/c'],
    [['/./a', 'b/c/.'], '/a/b/c'],
    [['//////', './.'], '/'],
    [['a/b', '../c'], 'a/c'],
    [['a/b', '../../../c'], '../c'],
    [['/a/b', '../../../..'], '/']
  ]

  tests.forEach(([input, output], index) => {
    cnt++
    try {
      const actual = joinPath(...input)
      strictEqual(actual, output, `    Expected '${output}' but get '${actual}'.`)
      console.log(`  ${index} \x1b[1;32mok\x1b[0m`)
    } catch (e) {
      ecnt++
      console.log(`  ${index} \x1b[1;31merror\x1b[0m`)
      console.log(e.message)
    }
  })
  console.log(`testJoinPath end: run ${cnt} tests, ` +
    (ecnt === 0
      ? `\x1b[1;32m${cnt - ecnt} ok\x1b[0m and no error.`
      : `${cnt - ecnt} ok and \x1b[1;31m${ecnt} error\x1b[0m.`))
}

const runMain = () => {
  testJoinPath()
}

runMain()
