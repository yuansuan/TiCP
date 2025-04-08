import * as React from 'react'
import Section from './Section'

interface InfoListProps {
  template: any
  dataSource: any
  transformer?: any
}

function InfoList(props: InfoListProps) {
  const { template, dataSource, transformer } = props

  return template.map((item, index) => (
    <Section
      key={index}
      title={item.title}
      icon={item.icon}
      type={item.type}
      partitionNum={item.partitionNum}
      children={item.children}
      dataSource={dataSource}
      transformer={transformer}
    />
  ))
}
export default InfoList
