import { InfoBlock, InfoCell } from '@/components'
import * as React from 'react'
import { InfoBlockTemplate } from './type'
import Residual from './Residual'

const _get = require('lodash.get')

interface SectionProps extends InfoBlockTemplate {
  dataSource: any
  transformer: any
}

function Section(props: SectionProps) {
  const { icon, title, type, partitionNum, dataSource, transformer } = props

  return (
    <InfoBlock title={title} icon={icon}>
      {type === 1 ? (
        <Residual />
      ) : (
        props.children.map((item, index) => {
          const share = item.share || 1
          const Value =
            dataSource && _get(dataSource, item.key)
              ? transformer && transformer[item.key]
                ? transformer[item.key](_get(dataSource, item.key))
                : _get(dataSource, item.key)
              : _get(dataSource, item.key) === 0
              ? 0
              : '--'

          return (
            <InfoCell
              key={index}
              infoKey={item.text}
              infoVal={
                Value && item.type && item.type === 1 ? (
                  <Value path={dataSource[item.key].split(':')[1]} />
                ) : (
                  Value
                )
              }
              infoValTip={item.type && item.type === 1 ? '' : Value}
              width={(share * 100) / partitionNum + '%'}
            />
          )
        })
      )}
    </InfoBlock>
  )
}
export default Section
