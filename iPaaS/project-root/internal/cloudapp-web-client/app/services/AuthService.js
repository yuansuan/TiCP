import Http from '@/services/AxiosHttp'

export default {
  verifyAuth(link) {
    return Http.get(`/visual/worktask/conn/${link}`)
  },
}
