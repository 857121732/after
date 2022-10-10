package robot

import (
	"errors"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Conf struct {
	Port        string `yaml:"port"`
	AppSecret   string `yaml:"app_secret"`
	CorporaPath string `yaml:"corpora_path"` // 语料库地址
	WebHook     string `yaml:"web_hook"`     // 钉钉机器人hook
}

var globalConf *Conf

func parseConf(file string) (*Conf, error) {
	confStr, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("Read conf file(%s) error(%v).\n", file, err)
		return nil, err
	}

	curConf := &Conf{}
	err = yaml.Unmarshal(confStr, curConf)
	if err != nil {
		log.Printf("Unmarshal conf(%s) error(%v).\n", string(confStr), err)
		return nil, err
	}

	if curConf.Port == "" || curConf.AppSecret == "" || curConf.CorporaPath == "" || curConf.WebHook == "" {
		log.Printf("Conf(%s) invalid.\n", string(confStr))
		return nil, errors.New("conf invalid")
	}

	globalConf = curConf
	return curConf, nil
}
