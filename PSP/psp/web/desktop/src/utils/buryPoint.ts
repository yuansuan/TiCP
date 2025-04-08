/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */
import { env } from '@/domain'

interface Props {
  category: string
  action?: string
  label?: string
  value?: number
}

export const buryPoint = ({
  category,
  action = '',
  label = '',
  value,
}: Props) => {
  ;(window['_czc'] || []).push([
    '_trackEvent',
    category,
    action,
    label || env.company?.id || 1,
    value,
  ])
}

// used for daovoice
window['__buryPoint__'] = buryPoint
