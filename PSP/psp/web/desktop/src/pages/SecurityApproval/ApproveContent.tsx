import React from 'react'

function Content({ title, children }) {
  return <span title={title}>{children}</span>
}

export function ApproveContent({ title, jsx }) {
  return (
    <Content title={title}>${jsx}</Content>
  )
}
