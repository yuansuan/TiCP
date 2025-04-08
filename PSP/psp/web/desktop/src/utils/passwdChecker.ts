enum RULES {
  'minLength' = 0,
  'maxLength',
  'Char_Lower',
  'Char_Upper',
  'Number',
  'Char_Special',
}

let passwdChecker: any = {
  configs: {
    strengthCheck: true,
    maxLength: 128,
    minLength: 10,
  },
  config(params) {
    Object.assign(this.configs, params)
  },
  rules: {},
}

passwdChecker.rules.tests = {
  required: [
    // enforce a minimum length
    password => {
      if (password.length < passwdChecker.configs.minLength) {
        return `密码长度不能小于${passwdChecker.configs.minLength}`
      }
      return null
    },

    // enforce a maximum length
    password => {
      if (password.length > passwdChecker.configs.maxLength) {
        return `密码长度不能大于${passwdChecker.configs.maxLength}`
      }
      return null
    },
  ],

  strengthCheck: [
    // require at least one lowercase letter
    password => {
      if (!/[a-z]/.test(password)) {
        return '密码至少包含一个小写字母'
      }
      return null
    },

    // require at least one uppercase letter
    password => {
      if (!/[A-Z]/.test(password)) {
        return '密码至少包含一个大写字母'
      }
      return null
    },

    // require at least one number
    password => {
      if (!/[0-9]/.test(password)) {
        return '密码至少包含一个数字'
      }
      return null
    },

    // require at least one special character
    password => {
      if (!/[^A-Za-z0-9]/.test(password)) {
        return '密码至少包含一个特殊字符'
      }
      return null
    },
  ],
}

// This method tests password
passwdChecker.test = function(password) {
  // create an object to store the test results
  let result = {
    errors: {},
    failedTests: [],
    passedTests: [],
    requiredTestErrors: {},
    strengthCheckTestErrors: {},
    strengthCheckTestsPassed: 0,
  }

  let i = 0
  this.rules.tests.required.forEach(function(test) {
    let err = test(password)

    let ruleName = RULES[i]

    if (typeof err === 'string') {
      result.errors[ruleName] = err
      result.requiredTestErrors[ruleName] = err
      result.failedTests.push(ruleName)
    } else {
      result.passedTests.push(ruleName)
    }
    i++
  })

  if (this.configs.strengthCheck) {
    let j = this.rules.tests.required.length

    this.rules.tests.strengthCheck.forEach(function(test) {
      let err = test(password)
      let ruleName = RULES[j]
      if (typeof err === 'string') {
        result.errors[ruleName] = err
        result.strengthCheckTestErrors[ruleName] = err
        result.failedTests.push(ruleName)
      } else {
        result.strengthCheckTestsPassed++
        result.passedTests.push(ruleName)
      }
      j++
    })
  }

  // return the result
  return result
}

export default passwdChecker
