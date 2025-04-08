import React, { useState, useEffect } from 'react'
import { Select } from 'antd'
import { currentUser } from '@/domain'
import { Http } from '@/utils'

const Option = Select.Option

export async function fetchProjectList(callback) {
  const res = await Http.get('/project/listForParam', {
    params: {
      is_admin: currentUser.hasSysMgrPerm
    }
  })
  callback(res?.data?.projects || [])
}

export function ProjectSelector({value, onSelect, onDeselect}) {
  const [projectList, setProjectList] = useState([])

  useEffect(() => {
    fetchProjectList((projects) => setProjectList(projects))
  }, [])

  return (
    <Select 
      style={{width: 200}}
      value={value} 
      maxTagCount={4}
      mode={'multiple'}
      onSelect={onSelect}
      onDeselect={onDeselect}>
      {
        projectList.map(p => 
          <Option 
            key={p.id} 
            value={p.id}> 
            {p.name} 
          </Option>
        )
      }
    </Select>  
  )
}