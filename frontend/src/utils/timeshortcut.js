/**
 * Today.
 * @returns {Array}
 */
export function today() {
  const end = new Date()
  const start = new Date()
  end.setTime(Math.ceil(+new Date() / 1000 / 3600 / 24) * 1000 * 3600 * 24 - 8 * 3600 * 1000)
  start.setTime(end.getTime() - 3600 * 1000 * 24 * 1)
  end.setTime(end.getTime() - 1000)
  return [start, end]
}

/**
 * Yesterday.
 * @returns {Array}
 */
export function yesterday() {
  const [a, b] = today()
  const start = new Date(a - 3600 * 1000 * 24)
  const end = new Date(b - 3600 * 1000 * 24)
  return [start, end]
}

/**
 * Day before yesterday.
 * @returns {Array}
 */
export function dbyesterday() {
  const [a, b] = yesterday()
  const start = new Date(a - 3600 * 1000 * 24)
  const end = new Date(b - 3600 * 1000 * 24)
  return [start, end]
}

/**
 * Lastday.
 * @param {Integer} days
 * @returns {Array}
*/
export function lastday(days) {
  const a = new Date()
  const b = +a
  a.setTime(b - 3600 * 1000 * 24 * days)
  return [a, new Date(b)]
}

/**
 * Lasthour.
 * @param {Integer} hours
 * @returns {Array}
*/
export function lasthour(hours) {
  const a = new Date()
  const b = +a
  a.setTime(b - 3600 * 1000 * hours)
  return [a, new Date(b)]
}
