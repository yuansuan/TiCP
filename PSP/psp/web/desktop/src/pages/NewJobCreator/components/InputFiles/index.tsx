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
import { Alert } from 'antd'
import { useStore } from '../../store'

interface IProps {
  fileTree: FileTree
}

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
        {/* {state.isFluent && (
          <Alert
            type='info'
            showIcon
            message={
              <StyledMessage>
                <div>Fluent求解需上传.jou的脚本文件，</div>
                <a
                  className='right'
                  href='./jou脚本编写指南.pdf'
                  download='jou脚本编写指南.pdf'>
                  下载jou脚本编写指南
                </a>
              </StyledMessage>
            }
            style={{ marginBottom: '10px' }}
          />
        )}
        {state.inputFileTip && (
          <Alert
            type='info'
            showIcon
            message={
              <StyledMessage
                dangerouslySetInnerHTML={{
                  __html: `请选择${state.inputFileTip
                    .map(
                      ([main, ...slaves]) =>
                        `<b>${main}</b>文件作为主文件${
                          slaves.length > 0
                            ? `，<b>${slaves.join('，')}</b>文件作为从文件`
                            : ''
                        }`
                    )
                    .join('<i>或</i>')}`
                }}
              />
            }
            style={{ marginBottom: '10px' }}
          />
        )} */}
        <FileList fileTree={fileTree} />
      </Provider>
    </InputFilesStyle>
  )
})
