module.exports = require('neostandard')({
  env: ['browser'],
  globals: {
    __DEBUG__: 'readonly',
    __PACKAGE_VERSION__: 'readonly'
  },
  files: ['src/**/*.js', 'scripts/**/*.js'],
  ignores: ['dist/**/*']
})
