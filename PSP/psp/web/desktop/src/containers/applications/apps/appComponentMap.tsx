/* Copyright (C) 2016-present, Yuansuan.cn */
import React from 'react'
import NewJobCreatorStarCCM from '@/pages/NewJobCreatorByApp'
import NewJobCreator from '@/pages/NewJobCreator'
import VisList from '@/pages/VisList/List'
import FileMGT from '@/pages/FileMGT'
import NewJobManager from '@/pages/NewJobManager'
import NewJobDetail from '@/pages/NewJobDetail'
import NewJobSetDetail from '@/pages/NewJobSetDetail'

export const AppComponentMap = {
  files: <FileMGT />,
  'new-job-creator': [<NewJobCreator />, <NewJobCreatorStarCCM />],
  'new-jobs': <NewJobManager />,
  'vis-session': <VisList />,
  'new-job': <NewJobDetail />,
  'new-job-set': <NewJobSetDetail />
}
