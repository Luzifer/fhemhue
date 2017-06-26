package main

import (
	"fmt"
	"io/ioutil"
	"net"

	yaml "gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pborges/huemulator"
	"github.com/satori/uuid"
)

type config struct {
	Telnet struct {
		IP   string `yaml:"ip"`
		Port int    `yaml:"port"`
	} `yaml:"telnet"`
	Devices []device `yaml:"devices"`
}

func loadConfig(filename string) (*config, error) {
	var err error
	if filename, err = homedir.Expand(filename); err != nil {
		return nil, err
	}

	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	res := &config{}
	return res, yaml.Unmarshal(d, res)
}

func (c config) Switch(uuid, state string) bool {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   net.ParseIP(c.Telnet.IP),
		Port: c.Telnet.Port,
	})
	if err != nil {
		log.Errorf("Unable to dial telnet connection: %s", err)
		return false
	}
	defer conn.Close()

	for _, d := range c.Devices {
		if d.UUID() != uuid {
			continue
		}

		if _, err := fmt.Fprintf(conn, "set %s %s\n", d.ID, d.States.GetStateCommand(state)); err != nil {
			log.Errorf("Unable to send command through telnet: %s", err)
			return false
		}
	}

	return true
}

func (c config) GetLights() []huemulator.Light {
	lights := []huemulator.Light{}

	for _, d := range c.Devices {
		lights = append(lights, huemulator.Light{
			UUID:    d.UUID(),
			Name:    d.Name,
			OnFunc:  func(l huemulator.Light) bool { return c.Switch(l.UUID, "on") },
			OffFunc: func(l huemulator.Light) bool { return c.Switch(l.UUID, "off") },
		})
	}

	return lights
}

type device struct {
	ID     string `yaml:"id"`
	Name   string `yaml:"name"`
	States states `yaml:"states"`
}

func (d device) UUID() string {
	ns := uuid.NewV5(uuid.NamespaceOID, "fhemhue")
	return uuid.NewV5(ns, d.ID).String()
}

type states map[string]string

func (s states) GetStateCommand(stateType string) string {
	if v, ok := s[stateType]; ok {
		return v
	}

	return stateType
}
