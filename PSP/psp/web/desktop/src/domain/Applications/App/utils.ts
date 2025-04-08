import {
  transform,
  isEqual,
  isObject,
  isArray,
  includes,
  find,
  differenceWith,
} from 'lodash'

export function arrayDifferenceByKey(newArr, base, key) {
  const result = { update: [], add: [], delete: [] };

  // get added items
  result.add = newArr.filter(item => {
    if (!isObject(item)) {
      return !base.includes(item);
    } else {
      return !base.some(baseItem => baseItem[key] === item[key]);
    }
  });

  // get deleted items
  result.delete = base.filter(item => {
    if (!isObject(item)) {
      return !newArr.includes(item);
    } else {
      return !newArr.some(newItem => newItem[key] === item[key]);
    }
  });

  // get updated items (only if key is same and value is different)
  result.update = base
    .map(item => {
      if (!isObject(item)) {
        return undefined;
      }

      const newItem = newArr.find(newItem => newItem[key] === item[key]);
      if (newItem && newItem.value !== item.value) {
        return { old: item, new: newItem };
      }
    })
    .filter(item => !!item);

  return result;
}



export function arrayDifference(newArr, base, key?) {
  const result = { update: [], add: [], delete: [] }

  // get added items
  result.add = newArr.filter(item => {
    if (!isObject(item)) {
      return !includes(base, item)
    } else {
      return !find(base, baseItem => baseItem[key] === item[key])
    }
  })

  // get deleted items
  result.delete = base.filter(item => {
    if (!isObject(item)) {
      return !includes(newArr, item)
    } else {
      return !find(newArr, newItem => newItem[key] === item[key])
    }
  })

  // get updated items(object)
  result.update = base
    .map(item => {
      if (!isObject(item)) {
        return undefined
      }

      const newItem = find(newArr, newItem => newItem[key] === item[key])
      if (!newItem) {
        return undefined
      }

      return { old: item, new: newItem }
    })
    .filter(item => !!item)

  return result
}

/**
 * Deep diff between two object, using lodash
 * @param  {Object} object Object compared
 * @param  {Object} base   Object to compare with
 * @return {Object}        Return a new object who represent the diff
 */
export function difference(object, base) {
  function changes(object, base) {
    const res = transform(object, function(result, value, key) {
      if (!isEqual(value, base[key])) {
        // ignore diff arrays
        if (isArray(value) && isArray(base[key])) {
          const isDiff =
            differenceWith(value, base[key]).length > 0 ||
            differenceWith(base[key], value).length > 0
          if (isDiff) {
            result[key] = []
          }
        } else if (isObject(value) && isObject(base[key])) {
          // recurse diff object
          result[key] = changes(value, base[key])
        } else {
          result[key] = { old: base[key], new: value }
        }
      }
    })
    return Object.keys(res).length > 0 ? res : null
  }
  return changes(object, base)
}

export const statusMap = {
  published: {
    text: '已发布',
    color: '#194E8B',
  },
  unpublished: {
    text: '未发布',
    color: '#FFB68F',
  },
}
