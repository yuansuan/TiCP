import { Http } from "@/utils"

export const monitorServer = {
  jobStateData(start_time: string, end_time: string) {
    return Http.get(`/dashboard/jobInfo?start=${start_time}&end=${end_time}`)
  },
  userJobData(start_time: string, end_time: string) {
    return Http.get(`/job/userJobNum?start=${start_time}&end=${end_time}`)
  },
  appJobData(start_time: string, end_time: string) {
    return Http.get(`/job/appJobNum?start=${start_time}&end=${end_time}`)
  },
  projectJobData(start_time: string, end_time: string) {
    return Http.get(`/job/statistics/top5ProjectInfo?start=${start_time}&end=${end_time}`)
  },
  resourceInfo(start_time: string, end_time: string) {
    return Http.get(`/dashboard/resourceInfo?start=${start_time}&end=${end_time}`)
  },
  clusterInfo() {
    return Http.get(`/dashboard/clusterInfo`)
  },
}