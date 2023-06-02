package logic_qq_wry

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/internal/logic"
	"github.com/zu1k/nali/pkg/qqwry"
	"github.com/zu1k/nali/pkg/wry"
	"github.com/zu1k/nali/pkg/zxipv6wry"
)

var (
	_ logic.QQWry = &QQWry{}
)

// QQWry
// @autowire(logic.QQWry,set=logics)
type QQWry struct {
	IPv4 *qqwry.QQwry
	IPv6 *zxipv6wry.ZXwry
}

func NewQQWry(config config.QQWry) (s *QQWry, err error) {
	var ipv4 *qqwry.QQwry
	if ipv4, err = loadIPv4(config); err != nil {
		return
	}
	var ipv6 *zxipv6wry.ZXwry
	if ipv6, err = loadIPv6(config); err != nil {
		return
	}
	s = &QQWry{IPv4: ipv4, IPv6: ipv6}
	return
}

func (s *QQWry) Find(ctx context.Context, query string) (location fmt.Stringer, err error) {
	_ = ctx
	query = strings.TrimSpace(query)
	ip := net.ParseIP(query)
	if ip == nil {
		err = errors.New("请输入正确的 IP 地址")
		return
	}
	if strings.Contains(query, ":") {
		return s.IPv6.Find(query)
	}
	return s.IPv4.Find(query)
}

func loadIPv4(config config.QQWry) (ipv4 *qqwry.QQwry, err error) {
	var data []byte
	if len(config.IPv4FilePath) <= 0 {
		data = config.IPv4Data
	} else {
		if _, err = os.Stat(config.IPv4FilePath); err != nil && os.IsNotExist(err) {
			return
		}
		var f *os.File
		if f, err = os.OpenFile(config.IPv4FilePath, os.O_RDONLY, 0400); err != nil {
			return
		}
		defer func() { _ = f.Close() }()
		if data, err = io.ReadAll(f); err != nil {
			return
		}
		if !qqwry.CheckFile(data) {
			err = errors.New("纯真 IPv4 库存在错误，请重新下载")
			return
		}
	}
	header := data[0:8]
	start := binary.LittleEndian.Uint32(header[:4])
	end := binary.LittleEndian.Uint32(header[4:])
	ipv4 = &qqwry.QQwry{IPDB: wry.IPDB[uint32]{Data: data, OffLen: 3, IPLen: 4, IPCnt: (end-start)/7 + 1, IdxStart: start, IdxEnd: end}}
	return
}

func loadIPv6(config config.QQWry) (ipv6 *zxipv6wry.ZXwry, err error) {
	var data []byte
	if len(config.IPv6FilePath) <= 0 {
		data = config.IPv6Data
	} else {
		if _, err = os.Stat(config.IPv6FilePath); err != nil && os.IsNotExist(err) {
			return
		}
		var f *os.File
		if f, err = os.OpenFile(config.IPv6FilePath, os.O_RDONLY, 0600); err != nil {
			return
		}
		defer func() { _ = f.Close() }()
		if data, err = io.ReadAll(f); err != nil {
			return
		}
		if !qqwry.CheckFile(data) {
			err = errors.New("纯真 IPv6 库存在错误，请重新下载")
			return
		}
	}
	header := data[:24]
	offLen := header[6]
	ipLen := header[7]
	start := binary.LittleEndian.Uint64(header[16:24])
	counts := binary.LittleEndian.Uint64(header[8:16])
	end := start + counts*11
	ipv6 = &zxipv6wry.ZXwry{IPDB: wry.IPDB[uint64]{Data: data, OffLen: offLen, IPLen: ipLen, IPCnt: counts, IdxStart: start, IdxEnd: end}}
	return
}
