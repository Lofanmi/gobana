
/**
 * This is just a simple version of deep copy
 * Has a lot of edge cases bug
 * If you want to use a perfect deep copy, use lodash's _.cloneDeep
 * @param {Object} source
 * @returns {Object}
 */
export function deepClone(source) {
  if (!source && typeof source !== 'object') {
    throw new Error('error arguments', 'deepClone')
  }
  const targetObj = source.constructor === Array ? [] : {}
  Object.keys(source).forEach(keys => {
    if (source[keys] && typeof source[keys] === 'object') {
      targetObj[keys] = deepClone(source[keys])
    } else {
      targetObj[keys] = source[keys]
    }
  })
  return targetObj
}

/**
 * date (adm_yy)
 * date('Y-m-d H:i:s')
 * @param {string} fmt
 * @param {Interger} timestamp
 */
export function date(fmt, timestamp) {
  fmt = fmt || 'Y-m-d H:i:s'
  let D = timestamp || +(new Date()) / 1000
  const week = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
  D = new Date((D >>> 0) * 1000)
  const dd = {
    year: D.getYear(),
    month: D.getMonth() + 1,
    date: D.getDate(),
    day: week[D.getDay()],
    hours: D.getHours(),
    minutes: D.getMinutes(),
    seconds: D.getSeconds()
  }
  dd.g = dd.hours > 12 ? Math.ceil(dd.hours / 2) : dd.hours
  const oType = {
    Y: D.getFullYear(),
    y: dd.year,
    m: dd.month < 10 ? '0' + dd.month : dd.month,
    n: dd.month,
    d: dd.date < 10 ? '0' + dd.date : dd.date,
    j: dd.date,
    D: dd.day,
    H: dd.hours < 10 ? '0' + dd.hours : dd.hours,
    h: dd.g < 10 ? '0' + dd.g : dd.g,
    G: dd.hours,
    g: dd.g,
    i: dd.minutes < 10 ? '0' + dd.minutes : dd.minutes,
    s: dd.seconds < 10 ? '0' + dd.seconds : dd.seconds
  }
  for (const i in oType) {
    fmt = ('' + fmt).replace(i, oType[i])
  }
  return fmt
}

/**
 * trim
 * @param {string} s
 */
export function trim(s) {
  s = s || ''
  return s.replace(/(^\s*)|(\s*$)/g, '')
}

/**
 * arrayChunk
 * @param {Array} arr
 * @param {Integer} size
 */
export function arrayChunk(arr, size) {
  const result = []
  for (let i = 0; i < arr.length; i = i + size) {
    result.push(arr.slice(i, i + size))
  }
  return result
}

/**
 * formatInteger
 * @param {Number} num 欲格式化的数字
 * @param {Number} n 保留的位数
 * @param {String} c 填充字符
 * @returns 格式化后的数字
 */
export function formatInteger(num, n, c = '0') {
  const numStr = String(num)
  var formatNum = (Array(n).join(c) + num).slice(-n)

  if (numStr.length > formatNum) {
    formatNum = numStr
  }

  return formatNum
}

/**
 * time33
 * @param {String} str 字符串
 * @returns 哈希值
 */
export function time33(str) {
  str = str || ''
  for (var i = 0, len = str.length, hash = 5381; i < len; ++i) {
    hash += (hash << 5) + str.charAt(i).charCodeAt()
  }
  return hash & 0x7fffffff
}
