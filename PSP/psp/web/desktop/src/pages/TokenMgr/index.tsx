import React, { useEffect } from 'react'
import { observer, useLocalStore } from 'mobx-react-lite'
import { List } from './List'
import { Page } from '@/components'
import styled from 'styled-components'
import { Button } from '@/components'
import { Http as v2Http } from '@/utils/v2Http'
import { env } from '@/domain'

export const StyledLayout = styled.div`
  .header {
    padding: 12px 20px;
    border-bottom: 1px solid ${({ theme }) => theme.borderColorBase};
  }

  .main {
    padding: 14px 20px;
  }
`

const ProjectMemberMGTPage = observer(function ProjectMemberMGTPage() {
  const store = useLocalStore(() => ({
    list: [],
    setTokenList(tokens) {
      this.list = tokens
    },
    loading: true,
    setLoading(flag) {
      this.loading = flag
    }
  }))

  const fetch = async () => {
    store.setLoading(true)
    const data = await v2Http.get(`/kms/token/${env.project?.id} `)
    store.setTokenList(data.list)
    store.setLoading(false)
  }

  useEffect(() => {
    fetch()
  }, [])

  const add = async () => {
    await v2Http.post(`/kms/apply/${env.project?.id}`)
    fetch()
  }

  return (
    <Page header={null}>
      <StyledLayout>
        <div className='header'>
          <Button type='primary' onClick={add}>
            新建密钥
          </Button>
        </div>
        <div className='main'>
          <List model={store} loading={store.loading} action={{ fetch }} />
        </div>
      </StyledLayout>
    </Page>
  )
})

export default ProjectMemberMGTPage
