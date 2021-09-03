/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"dubbo.apache.org/dubbo-go/v3/cluster/router/chain"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/router/chain"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"fmt"
	"github.com/creasty/defaults"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/pkg/errors"
)

// RouterConfig is the configuration of the router.
type RouterConfig struct {
	VirtualService  Property `yaml:"virtual_service" json:"virtual_service,omitempty" property:"virtual_service"`
	DestinationRule Property `yaml:"destination_rule" json:"destination_rule,omitempty" property:"destination_rule"`
}

type Property struct {
	Path  string `yaml:"path" json:"path,omitempty" property:"path"`
	Genre string `default:"yaml" yaml:"genre" json:"genre,omitempty" property:"genre"`
	Delim string `default:"." yaml:"delim" json:"delim,omitempty" property:"delim"`
}

// Prefix dubbo.router
func (RouterConfig) Prefix() string {
	return constant.RouterConfigPrefix
}

func (rc *RouterConfig) check() error {
	if err := defaults.Set(rc); err != nil {
		return err
	}
	return nil
}

func initRouterConfig(rc *RootConfig) error {
	router := rc.Router
	if router != nil {
		var err error
		if err = router.check(); err != nil {
			return err
		}

		var vsBytes, drBytes []byte
		vsBytes, err = getConfigBytes(router.VirtualService)
		if err != nil {
			return err
		}
		drBytes, err = getConfigBytes(router.DestinationRule)
		if err != nil {
			return err
		}
		chain.SetVSAndDRConfigByte(vsBytes, drBytes)
	}
	return nil
}

func getConfigBytes(p Property) ([]byte, error) {
	var (
		b   []byte
		err error
	)
	lc := NewLoaderConf(WithPath(p.Path), WithGenre(p.Genre), WithDelim(p.Delim))
	k := getKoanf(lc)
	switch p.Genre {
	case "yaml", "yml":
		b, err = k.Marshal(yaml.Parser())
	case "json":
		b, err = k.Marshal(json.Parser())
	case "toml":
		b, err = k.Marshal(toml.Parser())
	default:
		err = errors.New(fmt.Sprintf("Unsupported %s file type", p.Genre))
	}
	return b, err
}

//// LocalRouterRules defines the local router config structure
//type LocalRouterRules struct {
//	RouterRules []interface{} `yaml:"routerRules"`
//}
//
//// RouterInit Set config file to init router config
//func RouterInit(vsConfigPath, drConfigPath string) error {
//	vsBytes, err := yaml.LoadYMLConfig(vsConfigPath)
//	if err != nil {
//		return err
//	}
//	drBytes, err := yaml.LoadYMLConfig(drConfigPath)
//	if err != nil {
//		return err
//	}
//	chain.SetVSAndDRConfigByte(vsBytes, drBytes)
//	return nil
//}
