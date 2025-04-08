/* Copyright (C) 2016-present, Yuansuan.cn */

import { historyService } from './history'
import { showNewComerModal } from '@/pages/JobCreator/NewComer/NewComerModal'

export default {
  ...historyService,
  showNewComer: async () => {
    try {
      await showNewComerModal()
    } catch (err) {}
  }
}
