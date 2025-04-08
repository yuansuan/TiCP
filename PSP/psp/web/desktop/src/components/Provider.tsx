/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { observer } from 'mobx-react-lite'
import { ApolloProvider } from '@apollo/client'
import { apolloClient } from '@/utils'
import { ConfigProvider } from 'antd'
import { ConfigProviderProps } from 'antd/es/config-provider'
import { ThemeProvider } from 'styled-components'
import { theme } from '@/utils'
import zh_CN from 'antd/lib/locale-provider/zh_CN'

type Props = {
  children: React.ReactNode
  configProviderProps?: ConfigProviderProps
}

export const Provider = observer(function Provider({
  children,
  configProviderProps,
}: Props) {
  return (
    <ApolloProvider client={apolloClient}>
      <ThemeProvider theme={theme}>
        <ConfigProvider locale={zh_CN} {...configProviderProps}>
          {children}
        </ConfigProvider>
      </ThemeProvider>
    </ApolloProvider>
  )
})
