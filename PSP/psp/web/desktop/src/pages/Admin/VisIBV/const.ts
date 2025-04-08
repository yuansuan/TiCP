import { currentUser } from '@/domain'

export const hasPerm = currentUser.perms.includes('admin')
