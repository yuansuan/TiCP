import * as React from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react'
import flow from 'lodash/flow'

import { App } from '@/domain/Applications'
import { inject } from '@/pages/context'
import SectionUI from './Section'
import SectionTarget from './SectionTarget'

const Wrapper = styled.div`
  flex: 1;
  padding: 10px;
  overflow: auto;
  flex-direction: column;
  padding-bottom: 20px;
  position: relative;
  background-color: #f0f5fd;
`

interface IProps {
  app?: App
  formModel: any
}

@flow(observer, inject(({ app }) => ({ app })))
export default class Right extends React.Component<IProps> {
  render() {
    const {
      formModel,
      app: {
        subForm,
        subForm: { sections }
      }
    } = this.props

    return (
      <Wrapper>
        {sections.map((item, index) => (
          <SectionUI
            key={item._key}
            formModel={formModel}
            model={item}
            index={index}
            onDelete={subForm.delete}
          />
        ))}
        <SectionTarget index={sections.length} />
      </Wrapper>
    )
  }
}
