/* Copyright (C) 2016-present, Yuansuan.cn */

import { apolloClient } from '@/utils'
import { gql } from '@apollo/client'

export const JOBCOUNT_MONITOR = gql`
  query jobCountMonitor($payload: ChartInput!) {
    jobCountMonitor(payload: $payload) {
      date
      type
      value
    }
  }
`

export const COREHOUR_MONITOR = gql`
  query coreHourMonitor($payload: ChartInput!) {
    coreHourMonitor(payload: $payload) {
      date
      type
      value
    }
  }
`

const statsFragment = gql`
  fragment statsFragment on ChartData {
    type
    value
    date
  }
`

export const COREHOUR_STATS = gql`
  query($payload: ChartInput!) {
    coreHourStats(payload: $payload) {
      ...statsFragment
    }
  }
  ${statsFragment}
`

export const APP_COREHOUR_STATS = gql`
  query($payload: ChartInput!) {
    appCoreHourStats(payload: $payload) {
      ...statsFragment
    }
  }
  ${statsFragment}
`

export const JOBCOUNT_STATS = gql`
  query($payload: ChartInput!) {
    jobCountStats(payload: $payload) {
      ...statsFragment
    }
  }
  ${statsFragment}
`

export const PROJECT_OVERVIEW = gql`
  query {
    projectOverview {
      members
      amount
      jobs
      coreHours
    }
  }
`

export const dashboardServer = {
  async getProjectOverview() {
    return await apolloClient.query({
      query: PROJECT_OVERVIEW,
      fetchPolicy: 'network-only'
    })
  }
}
