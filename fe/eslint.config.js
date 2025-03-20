import neostandard, { resolveIgnoresFromGitignore } from 'neostandard'
import globals from 'globals'

export default [
  ...neostandard({
    env: ['browser'],
    globals: {
      __DEBUG__: 'readonly',
      __PACKAGE_VERSION__: 'readonly'
    },
    files: ['src/**/*.js', 'scripts/**/*.js'],
    ignores: resolveIgnoresFromGitignore()
  }),
  {
    languageOptions: { ecmaVersion: 2025, globals: globals.es2025 }
  }
]
