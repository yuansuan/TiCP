import * as React from 'react'
import styled from 'styled-components'

import { InlineDiff } from './Diff'

interface IProps {
  fieldDiff: any
}

const StyledOptions = styled.div`
  display: flex;
`

export default class FieldDiff extends React.Component<IProps> {
  render() {
    const { fieldDiff } = this.props

    return (
      <div>
        <div>
          <span className='tag'>更新组件（{fieldDiff.key}）</span>
        </div>
        {Object.keys(fieldDiff.props).map(key => {
          const fieldItem = fieldDiff.props[key]
          if (key === 'options' || key === 'default_values') {
            const showAdd = fieldItem.add.length > 0
            const showDelete = fieldItem.delete.length > 0

            return (
              <StyledOptions key={key}>
                {(showAdd || showDelete) && (
                  <span className='label'>{key}</span>
                )}
                {showDelete &&
                  fieldItem.delete.map(name => (
                    <InlineDiff key={name} Old={name} />
                  ))}
                {showAdd &&
                  fieldItem.add.map(name => (
                    <InlineDiff key={name} New={name} />
                  ))}
              </StyledOptions>
            )
          }

          return (
            <InlineDiff
              key={key}
              name={key}
              Old={fieldItem.old}
              New={fieldItem.new}
            />
          )
        })}
      </div>
    )
  }
}
