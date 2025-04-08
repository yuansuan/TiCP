/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { InputFilesStyle } from './style'
import { FileList } from './FileList'
import { Toolbar } from './Toolbar'
import { observer, useLocalStore } from 'mobx-react-lite'
import { FileTree } from '@/domain/JobBuilder/FileTree'
import { Provider } from './store'
import { useStore } from '../../store'

interface IProps {
  fileTree: FileTree
}

const StyledMessage = styled.div`
  display: flex;

  > b {
    color: ${({ theme }) => theme.primaryColor};
    margin: 0 4px;
  }

  > i {
    margin: 0 4px;
    font-weight: 400;
  }

  > .right {
    color: ${({ theme }) => theme.linkColor};
    text-decoration: none;
    cursor: pointer;
  }
`

export const InputFiles = observer(function InputFiles({ fileTree }: IProps) {
  const store = useStore()
  const state = useLocalStore(() => ({
    get type() {
      const type = store.data?.currentApp?.type
      return type && type.toLowerCase()
    },
    get isFluent() {
      return this.type?.includes('fluent')
    },
    get inputFileTip() {
      return false
      const tip = Object.entries({
        Abaqus: 'inp',
        'STAR-CCM+': 'sim',
        Optistruct: 'fem',
        'MSC Nastran': 'dat',
        'Ls-Dyna': 'k',
        Radioss: 'bdf',
        COMSOL: 'mph',
        Workbench: [['wbpj', '对应的文件夹'], ['wbpz']],
        'MSC Marc': 'dat',
        'MSC Adams': 'acf',
        'MSC scSTREAM': 's',
        'MSC SCTpre': 'j',
        'MSC SCTso': [['s', 'pre']],
        'MSC scFlowso': [['sph', 'gph']],
        'MSC Dytran': 's',
        'ANSYS Maxwell': 'aedt',
        'ANSYS HFSS': [['aedt (HFSS)'], ['aedt (HFSS Layout)', '同名文件夹']],
        Fluent: 'jou',
        'Code Aster': [['comm', 'export', 'mmed']],
        Slwave: 'siw',
        Icepak: 'tzr',
        Converge: 'inputs.in',
        Telemac: 'cas',
        FEKO: 'fdk',
        CFX: 'def',
        // Mechanical: 'inp',
        SW: 'in',
        OpenFOAM: 'sh'
      }).find(([type]) => type.toLowerCase() === this.type)?.[1]

      if (!tip) {
        return ''
      }

      return typeof tip === 'string' ? [[tip]] : tip
    }
  }))
  return (
    <InputFilesStyle>
      <Provider>
        <Toolbar fileTree={fileTree} />

        <FileList fileTree={fileTree} />
      </Provider>
    </InputFilesStyle>
  )
})
