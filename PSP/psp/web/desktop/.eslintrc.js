/* Copyright (C) 2016-present, Yuansuan.cn */

module.exports = {
  env: {
    browser: true,
    es6: true
  },
  plugins: ['header', 'eslint-plugin-react', '@typescript-eslint'],
  extends: [
    'plugin:react-hooks/recommended',
    'prettier',
    'prettier/@typescript-eslint'
  ],
  parser: '@typescript-eslint/parser',
  rules: {
    'header/header': [
      0,
      'block',
      [
        {
          pattern: ' Copyright',
          template: ' Copyright (C) 2016-present, Yuansuan.cn '
        }
      ]
    ],
    '@typescript-eslint/member-delimiter-style': [
      'error',
      {
        multiline: {
          delimiter: 'none',
          requireLast: true
        },
        singleline: {
          delimiter: 'semi',
          requireLast: false
        }
      }
    ],
    '@typescript-eslint/prefer-namespace-keyword': 'error',
    '@typescript-eslint/quotes': ['error', 'single'],
    '@typescript-eslint/semi': ['error', 'never'],
    '@typescript-eslint/type-annotation-spacing': 'error',
    'arrow-parens': ['off', 'always'],
    'brace-style': ['error', '1tbs'],
    'comma-dangle': 'off',
    'eol-last': 'off',
    'id-blacklist': [
      'error',
      'any',
      'Number',
      'number',
      'String',
      'string',
      'Boolean',
      'boolean',
      'Undefined',
      'undefined'
    ],
    'id-match': 'error',
    'linebreak-style': 'off',
    'max-len': 'off',
    'new-parens': 'off',
    'newline-per-chained-call': 'off',
    'no-eval': 'error',
    'no-extra-semi': 'off',
    'no-irregular-whitespace': 'off',
    'no-multiple-empty-lines': 'off',
    'no-shadow': [
      'off',
      {
        hoist: 'all'
      }
    ],
    'no-trailing-spaces': 'error',
    'no-unsafe-finally': 'error',
    'no-var': 'error',
    'quote-props': 'off',
    'react/jsx-curly-spacing': 'off',
    'react/jsx-equals-spacing': 'off',
    'react/jsx-wrap-multilines': 'off',
    'space-before-function-paren': 'off',
    'space-in-parens': ['off', 'never']
  }
}
