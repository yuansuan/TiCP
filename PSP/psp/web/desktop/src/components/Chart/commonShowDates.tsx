import React from 'react'
import moment from 'moment'

export const date = (showDates: any) => {
  return (
    <div style={{ height: 25, textAlign: 'center' }}>
      {moment(showDates[0]).format('YYYY-MM-DD HH:mm:ss') +
        ' ~ ' +
        moment(showDates[1]).format('YYYY-MM-DD HH:mm:ss')}
    </div>
  )
}
