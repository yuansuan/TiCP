import * as React from 'react'

import { InlineDiff } from './Diff'
import FieldDiff from './FieldDiff'
import FormItem from './FormItem'
import Field from '@/domain/Applications/App/Field'

interface IProps {
  sectionDiff: any
}

export default class SectionDiff extends React.Component<IProps> {
  render() {
    const { sectionDiff } = this.props
    const showAdd = sectionDiff.add.length > 0
    const showDelete = sectionDiff.delete.length > 0
    const showUpdate = sectionDiff.update.length > 0

    return (
      <div style={{ paddingLeft: 20 }}>
        {showDelete && (
          <div style={{ display: 'flex' }}>
            {sectionDiff.delete.map(item => (
              <InlineDiff
                key={item.id}
                Old={<FormItem model={new Field(item)} />}
              />
            ))}
          </div>
        )}

        {showAdd && (
          <div style={{ display: 'flex' }}>
            {sectionDiff.add.map(item => (
              <InlineDiff
                key={item.id}
                New={<FormItem model={new Field(item)} />}
              />
            ))}
          </div>
        )}

        {showUpdate && (
          <>
            {sectionDiff.update.map(item => (
              <FieldDiff key={item.key} fieldDiff={item} />
            ))}
          </>
        )}
      </div>
    )
  }
}
