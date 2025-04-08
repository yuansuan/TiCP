import { Http } from '@/utils'
export default class Auth {
  static async login(username, password) {
    if (username === '') {
      return Promise.reject()
    }
    if (password === '') {
      return Promise.reject()
    }

    return await Http.post('/auth/login', {
      name: username,
      password
    }).then(res => {
      return res
    })
  }

  static getSalt(username: string) {
    return Http.get('/auth/salt', {
      params: { username }
    }).then(res => res.data)
  }
}
