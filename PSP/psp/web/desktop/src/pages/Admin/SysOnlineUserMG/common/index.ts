//common sorting method sorted by session start time
export function sortByStartTime(sortData, sortType) {
  sortData = sortData.slice().sort((x, y) => {
    const xTime = new Date(x.created_at).getTime()
    const yTime = new Date(y.created_at).getTime()
    return sortType === 'asc' ? xTime - yTime : yTime - xTime
  })
  return sortData
}
