module.exports = {
  env: {
    es6: true,
    browser: true
  },
  extends: ['plugin:react/recommended', 'plugin:react-hooks/recommended', 'standard-with-typescript', 'standard-react', 'standard-jsx'],
  plugins: ['react-hooks', '@typescript-eslint'],
  ignorePatterns: ['.eslintrc.cjs', 'dist', '.yarn/*', '.pnp.*'],
  parser: '@typescript-eslint/parser',
  parserOptions: {
    project: require('path').join(__dirname, 'tsconfig.json')
  },
  rules: {
    'react/react-in-jsx-scope': 'off'
  }
}
