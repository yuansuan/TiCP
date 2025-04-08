import { Http } from '@/utils'

export type FilterParams = Partial<{
  file_name: string
  key: string
  start_seconds: string
  end_seconds: string
}>

export const fileRecordServer = {
  getList: querys =>
    Http.get('/filerecord/list', {
      params: {
        ...querys
      }
    })
}
