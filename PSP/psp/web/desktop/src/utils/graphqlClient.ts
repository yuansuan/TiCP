/* Copyright (C) 2016-present, Yuansuan.cn */

import {
    ApolloClient,
    ApolloClientOptions,
    HttpLink,
    ApolloLink,
  } from '@apollo/client'
  import { onError } from '@apollo/client/link/error'
  import message from 'antd/lib/message'
  import 'antd/lib/message/style/index.css'
  import { single } from './single'
  
  
  export const createGraphqlClient = function <TCacheShape>(
    apolloOptions: Omit<ApolloClientOptions<TCacheShape>, 'link'> & {
      link?: ApolloLink | ((link: ApolloLink) => ApolloLink)
    }
  ) {
    const httpLink = new HttpLink({ uri: apolloOptions.uri })
    const errorHandlerLink = onError(
      ({ networkError, graphQLErrors, operation }) => {
        // handle network error
        if (
          networkError?.name === 'TypeError' &&
          networkError?.message === 'Failed to fetch'
        ) {
          single(
            'http-client-network-error-message',
            () =>
              new Promise((resolve, reject) =>
                message.error('网络异常').then(resolve, reject)
              )
          )
          return
        }
  
        const error = (graphQLErrors || [])[0]
        if (!error) {
          if (networkError?.message) {
            single(
              'http-client-server-parse-error-message',
              () =>
                new Promise((resolve, reject) =>
                  message.error(networkError?.message).then(resolve, reject)
                )
            )
          }
  
          return
        }
  
        const context = operation.getContext()
        const disableErrorMessage = !!context?.disableErrorMessage
        const formatErrorMessage = context?.formatErrorMessage || (msg => msg)
        const exception = error?.extensions?.exception
  
        // handle business error
        if (error?.extensions?.code === 'INTERNAL_SERVER_ERROR') {
          const res = exception?.response
          const success = res?.success
          const msg = res?.message || error.message
  
          if (!success && !disableErrorMessage) {
            message.error(formatErrorMessage(msg, exception))
          }
        } else {
          // handle graphql error
          if (!disableErrorMessage) {
            message.error(formatErrorMessage(error.message, exception))
          }
        }
      }
    )
  
    let finalLink = errorHandlerLink.concat(httpLink)
    // ovrride default link
    if (apolloOptions?.link) {
      if (apolloOptions?.link instanceof ApolloLink) {
        finalLink = apolloOptions?.link
      } else {
        finalLink = apolloOptions?.link(finalLink)
      }
    }
  
    return new ApolloClient({
      ...apolloOptions,
      link: finalLink,
    })
  }
  

type GqlException = {
  status: number
  response: {
    errorCode: number
    success: boolean
    message: string
  }
}

export const createHeaderMiddleware = (headers: object | (() => object)) =>
  new ApolloLink((operation, forward) => {
    operation.setContext(({ headers: originHeaders }) => ({
      headers: {
        ...originHeaders,
        ...(typeof headers === 'function' ? headers() : headers),
      },
    }))

    return forward(operation)
  })

export const createErrorMiddleware = (
  handler: (exception: GqlException) => void
) =>
  onError(({ graphQLErrors }) => {
    const error = (graphQLErrors || [])[0]
    if (!error) {
      return
    }

    const exception = error?.extensions?.exception
    handler(exception)
  })
