import React from 'react'

interface IProps {
  selectedName: string[]
}
const SelectedItems = ({ selectedName }: IProps) => {
  return (
    <div style={{ maxHeight: '500px', overflowY: 'auto' }}>
      {selectedName.length === 1 ? (
        `确认要删除文件 ${selectedName[0]} 吗？`
      ) : (
        <div>
          <p>确认要删除如下文件吗？</p>
          <ul style={{ marginLeft: 20 }}>
            {selectedName.map((item, index) => (
              <li key={index}>{item}</li>
            ))}
          </ul>
        </div>
      )}
    </div>
  )
}

export default SelectedItems
