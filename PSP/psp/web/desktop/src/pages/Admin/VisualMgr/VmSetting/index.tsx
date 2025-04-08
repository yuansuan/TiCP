import React, { useEffect } from 'react'
import { observer } from 'mobx-react'

import { Wrapper } from './style'
import AllConfig from './AllConfig'
import currentUser from '@/domain/User'
import { history } from '@/utils'

const VirtualMachineConfig = observer(() => {
  useEffect(() => {
    const hasVisitedPerm = currentUser.perms.includes('system-system_config')
    if (!hasVisitedPerm) {
      history.push('/login')
    }
  }, [])

  return (
    <Wrapper>
      <div className='body'>
        <AllConfig />
      </div>
    </Wrapper>
  )
})

export default VirtualMachineConfig
