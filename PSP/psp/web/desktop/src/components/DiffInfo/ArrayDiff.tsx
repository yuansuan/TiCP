import React from 'react'
import { InlineDiff } from './Diff'

import { StyledSection } from './style'

interface IProps {
  arrayDiff: any
  title: string
}

export default class ArrayDiff extends React.Component<IProps> {
  render() {
    const { arrayDiff,title } = this.props
    const showAdd = arrayDiff.add.length > 0
    const showDelete = arrayDiff.delete.length > 0
    const showUpdate = arrayDiff.update.length > 0

    if (!showAdd && !showDelete && !showUpdate) {
      return null
    }

    return (
      <StyledSection>
        {(showAdd || showDelete || showUpdate) && (
          <div className='tag'>{title}</div>
        )}
        {showDelete && (
          <>
            {arrayDiff.delete.map(item => (
              <InlineDiff key={item} Old={<span>{item}</span>} />
            ))}
          </>
        )}

        {showAdd && (
          <>
            {arrayDiff.add.map(item => (
              <InlineDiff key={item} New={<span>{item}</span>} />
            ))}
          </>
        )}

        {showUpdate && (
          <>
            {arrayDiff.update.map(item => (
              <div key={item}>
                <InlineDiff key={item} Old={item} New={item} />
              </div>
            ))}
          </>
        )}
      </StyledSection>
    )
  }
}
