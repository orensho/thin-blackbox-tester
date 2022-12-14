package config

import (
	"fmt"
	"sync"
)

var lock = &sync.Mutex{}

type Proxy struct {
	Hostname string
	Customer string
}

type ProxiesSettings struct {
	Proxies []Proxy
}

var singleInstance *ProxiesSettings

func GetProxiesSettingsInstance() *ProxiesSettings {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			fmt.Println("Creating single instance now.")
			singleInstance = &ProxiesSettings{}
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return singleInstance
}

func (s *ProxiesSettings) Evaluate() error {
	return nil
}

func (s *ProxiesSettings) Validate() error {
	return nil
}

func (s *ProxiesSettings) SetProxies(proxies []Proxy) {
	s.Proxies = proxies
}

func (s *ProxiesSettings) GetProxies() *[]Proxy {
	return &s.Proxies
}
