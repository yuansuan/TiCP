import { CodeEditor } from '@/components'
import { observer } from 'mobx-react'
import * as React from 'react'

import { inject } from '@/pages/context'
import { App } from '@/domain/Applications'

interface IProps {
  app?: App
}

@inject(({ app }) => ({ app }))
@observer
export default class Script extends React.Component<IProps> {
  public editorRef = null

  public getValue() {
    return this.editorRef.getValue()
  }

  public render() {
    const { scriptData } = this.props.app

    return (
      <div
        style={{ height: '100%', backgroundColor: 'white', padding: '20px 0' }}>
        <CodeEditor
          ref={ref => (this.editorRef = ref)}
          value={scriptData}
          language='shell'
          readOnly={true}
        />
      </div>
    )
  }
}
