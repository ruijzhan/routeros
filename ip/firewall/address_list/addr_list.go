package addresslist

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ruijzhan/routeros"
)

type Entry struct {
	ID       string `json:"id,omitempty"`
	List     string `json:"list,omitempty"`
	Address  string `json:"address,omitempty"`
	Comment  string `json:"comment,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
	Timeout  string `json:"timeout,omitempty"`
}

func (e *Entry) String() string {
	bytes, _ := json.Marshal(e)
	return string(bytes)
}

func WithListName(list string) routeros.ListOptions {
	return func(cmd string) string {
		return fmt.Sprintf("%s ?list=%s", cmd, list)
	}
}

func List(cli *routeros.Client, opts ...routeros.ListOptions) ([]*Entry, error) {
	cmd := "/ip/firewall/address-list/print"
	for _, opt := range opts {
		cmd = opt(cmd)
	}

	reply, err := cli.RunArgs(strings.Split(cmd, " "))
	if err != nil {
		return nil, fmt.Errorf("list address-list error: %w", err)
	}

	list := make([]*Entry, len(reply.Re))
	for i, r := range reply.Re {
		m := r.Map
		disabled, _ := strconv.ParseBool(m["disabled"])
		list[i] = &Entry{
			ID:       m[".id"],
			List:     m["list"],
			Address:  m["address"],
			Comment:  m["comment"],
			Disabled: disabled,
			Timeout:  m["timeout"],
		}
	}
	return list, nil
}

func Add(cli *routeros.Client, list, address, timeout, comment string) error {
	cmd := fmt.Sprintf("/ip/firewall/address-list/add =list=%s =address=%s =timeout=%s =comment=%s",
		list, address, timeout, comment)
	_, err := cli.RunArgs(strings.Split(cmd, " "))
	return err
}
