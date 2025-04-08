export const propertyMapReduce = (hardward: Array<any>, property: string) => {
  const result = {}
  hardward.map(item => {
    if (!result[item[property]]) {
      result[item[property]] = true
    }
  })
  return Object.keys(result)
}
