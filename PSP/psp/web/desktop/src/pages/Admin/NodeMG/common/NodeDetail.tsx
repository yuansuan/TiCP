import { observer } from 'mobx-react'
import * as React from 'react'
import { InfoBlock, InfoCell } from '@/components'
import { NODE_INFO_CONF } from './const'

@observer
export default class NodeDetail extends React.Component<any> {
  render() {
    const { node } = this.props

    return (
      <div style={{ marginTop: 20 }}>
        {node &&
          NODE_INFO_CONF.map((config, index) => {
            const { children, partitionNum } = config

            return (
              <InfoBlock key={index} title={config.title}>
                {children.map((item, i) => {
                  const values = node
                  const { key, text, keyTip, formatter } = item
                  let value = values
                    ? values[key] !== undefined && values[key] !== ''
                      ? values[key]
                      : '--'
                    : '--'
                  value = formatter ? formatter(values) : value

                  return (
                    <InfoCell
                      key={i}
                      infoKey={text}
                      infoVal={value}
                      infoKeyTip={keyTip}
                      infoValTip={value === '--' ? null : value}
                      width={100 / partitionNum + '%'}
                    />
                  )
                })}
              </InfoBlock>
            )
          })}
      </div>
    )
  }
}
