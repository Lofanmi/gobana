
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
 * @param {Function} func
 * @param {number} wait
 * @param {boolean} immediate
 * @return {*}
 */
export function debounce(func, wait, immediate) {
  let timeout, args, context, timestamp, result

  const later = function() {
    // 据上一次触发时间间隔
    const last = +new Date() - timestamp

    // 上次被包装函数被调用时间间隔 last 小于设定时间间隔 wait
    if (last < wait && last > 0) {
      timeout = setTimeout(later, wait - last)
    } else {
      timeout = null
      // 如果设定为immediate===true，因为开始边界已经调用过了此处无需调用
      if (!immediate) {
        result = func.apply(context, args)
        if (!timeout) context = args = null
      }
    }
  }

  return function(...args) {
    context = this
    timestamp = +new Date()
    const callNow = immediate && !timeout
    // 如果延时不存在，重新设定延时
    if (!timeout) timeout = setTimeout(later, wait)
    if (callNow) {
      result = func.apply(context, args)
      context = args = null
    }

    return result
  }
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

/**
 * defaultInterval
 * @param {Number} timeA 开始时间
 * @param {Number} timeB 结束时间
 * @returns 区间范围/秒
 */
export function defaultInterval(timeA, timeB) {
  const maxChartPoints = 60
  let interval = parseInt((timeB - timeA) * 0.001 / maxChartPoints)
  if (interval <= 1) {
    interval = 1
  } else if (interval <= 5) {
    interval = 5
  } else if (interval <= 10) {
    interval = 10
  } else if (interval <= 30) {
    interval = 30
  } else if (interval <= 60) {
    interval = 60
  } else if (interval <= 300) {
    interval = 300
  } else if (interval <= 900) {
    interval = 900
  } else if (interval <= 1800) {
    interval = 1800
  } else if (interval <= 3600) {
    interval = 3600
  } else if (interval <= 3600 * 3) {
    interval = 3600 * 3
  } else if (interval <= 3600 * 9) {
    interval = 3600 * 9
  } else if (interval <= 3600 * 12) {
    interval = 3600 * 12
  } else {
    interval = 3600 * 24
  }
  return interval
}
