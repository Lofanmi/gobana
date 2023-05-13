package logic_lua_state

import (
	"unsafe"

	"go.uber.org/zap/buffer"
)

var (
	p buffer.Pool
)

func init() {
	p = buffer.NewPool()
}

func string2Bytes(s string) []byte {
	if s == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// nginxDecode nginx日志解析
func nginxDecode(s string) string {
	if s[0] != '{' {
		return ""
	}
	state := 0
	buf := p.Get()
	data := string2Bytes(s)
	length := len(data)
	for i := 0; i < length; i++ {
		letter := data[i]
		switch letter {
		case '\\':
			if state == 0 {
				state = 1
			}
		case 'x':
			if state == 1 {
				state = 2
			} else {
				buf.AppendByte('x')
			}
		default:
			if state == 2 {
				a, b := data[i], data[i+1]
				i++
				if a >= '0' && a <= '9' {
					a -= '0'
				} else if a >= 'A' && a <= 'F' {
					a -= 'A'
					a += 10
				}
				a *= 16
				if b >= '0' && b <= '9' {
					b -= '0'
				} else if b >= 'A' && b <= 'F' {
					b -= 'A'
					b += 10
				}
				a += b
				if a == '"' {
					buf.AppendByte('\\')
					buf.AppendByte('"')
				} else if a == '\\' {
					buf.AppendByte('\\')
					buf.AppendByte('\\')
				} else if a == '\t' {
					buf.AppendByte('\\')
					buf.AppendByte('t')
				} else if a == '\r' {
					buf.AppendByte('\\')
					buf.AppendByte('r')
				} else if a == '\n' {
					buf.AppendByte('\\')
					buf.AppendByte('n')
				} else {
					buf.AppendByte(a)
				}
			} else {
				buf.AppendByte(data[i])
			}
			state = 0
		}
	}
	s = buf.String()
	buf.Free()
	return s
}
