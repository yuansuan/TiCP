import { history } from '@/utils'
import { currentUser } from '@/domain'

export interface ListQuery {
  page: number
  pageSize: number
  query: string
}

export const defaultListQuery: ListQuery = {
  page: 1,
  pageSize: 10,
  query: ''
}

export async function checkPerm() {
  await currentUser.fetch()

  // if (currentUser.perms?.includes('system-user_management')) {
  //   return
  // }
  // history.replace('/sys/user')
}

export function toExpiredDate(expired_date) {
  return expired_date && expired_date.formatDate('YYYY-MM-DD 00:00:00')
    ? expired_date.formatDate('YYYY-MM-DD 00:00:00')
    : '永不过期'
}
