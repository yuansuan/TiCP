/* Copyright (C) 2016-present, Yuansuan.cn */

import { BaseVisualConfig } from '@/domain/Visualization/VisualConfig'

export const initialValues: BaseVisualConfig = {
  isOpen: false,
  activeTerminal: 0,
  maxTerminal: 0,
  bundleUsages: [],
}

export const visualServer = {
  fetch: async () => ({ ...initialValues }),
}
