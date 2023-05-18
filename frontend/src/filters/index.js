
/**
 * 格式化 JSON 字符串
 * @param {String} string
 */
export function jsonFormatter(string) {
  if (
    (string.charAt(0) === '{' && string.charAt(string.length - 1) === '}') ||
    (string.charAt(0) === '[' && string.charAt(string.length - 1) === ']')
  ) {
    return JSON.stringify(JSON.parse(string), null, 2)
  }
  return string
}

/**
 * 默认值
 * @param {String} string
 */
export function def(string) {
  return string || '-'
}
