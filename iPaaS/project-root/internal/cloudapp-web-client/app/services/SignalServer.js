import Http from '@/services/AxiosHttp'

export default {
  getSignalAddress() {
    return Http.get(`/signal/address`)
  }
}
