package traefik_plugin_ip2location_redirect

import (
	"context"
	"fmt"
	"net"
	"net/http"
)


type Config struct {
	// ip2location数据库文件路径
	Filename string `json:"filename"`
	// 地区缩写，多个地区以英文逗号分隔，如 CN,TW,HK,US,UK
	Regions []string `json:"regions"`
	// 跳转的url地址
	RedirectUrl string `json:"redirectUrl"`
	// 匹配反转，为true表示未匹配上地区时跳转，为false表示匹配上地区时跳转。默认为否
	NoMatch bool `json:"noMatch,omitempty"`
	// 是否永久重定向，默认为否
	Permanent bool `json:"permanent,omitempty"`
	// 获取IP的头信息字段名称，默认为remote-addr
	FromHeader string `json:"fromHeader,omitempty"`
	// 是否禁用出错时头信息，默认false
	DisableErrorHeader bool `json:"disableErrorHeader,omitempty"`
}


func CreateConfig() *Config {
	return &Config{}
}

type IP2LocationRedirect struct {
	next http.Handler
	name string
	config *Config
	db *DB
}

func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	db, err := OpenDB(config.Filename)
	if err != nil {
		return nil, fmt.Errorf("error open database file, %w", err)
	}

	return &IP2LocationRedirect{
		next: next,
		name: name,
		config: config,
		db: db,
	}, nil
}


func (a *IP2LocationRedirect) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ip, err := a.getIP(req)
	if err != nil {
		if !a.config.DisableErrorHeader {
			req.Header.Add("X-IP2LOCATION-REDIRECT-ERROR", err.Error())
		}
		a.next.ServeHTTP(rw,req)
		return
	}

	record, err := a.db.Get_all(ip.String())
	if err != nil {
		if !a.config.DisableErrorHeader {
			req.Header.Add("X-IP2LOCATION-REDIRECT-ERROR", err.Error())
		}
		a.next.ServeHTTP(rw,req)
		return
	}

	redirect := false
	if a.config.NoMatch {
		if !InArray(record.Country_short, a.config.Regions) {
			redirect = true
		}
	} else {
		if InArray(record.Country_short, a.config.Regions) {
			redirect = true
		}
	}

	if redirect {
		rw.Header().Set("Location", a.config.RedirectUrl)
		if a.config.Permanent {
			rw.WriteHeader(http.StatusMovedPermanently)
		} else {
			rw.WriteHeader(http.StatusFound)
		}
		return
	}
	a.next.ServeHTTP(rw, req)
}

func (a *IP2LocationRedirect) getIP(req *http.Request) (net.IP, error) {
	if a.config.FromHeader != "" {
		return net.ParseIP(req.Header.Get(a.config.FromHeader)), nil
	}

	addr, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, err
	}

	return net.ParseIP(addr), nil
}