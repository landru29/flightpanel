package detect

import (
	"fmt"
	"regexp"

	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/registry"
)

func listSerial(log logrus.FieldLogger, pattern *regexp.Regexp) ([]string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `HARDWARE\\DEVICEMAP\\SERIALCOMM`, registry.QUERY_VALUE)
	if err != nil {
		return nil, fmt.Errorf("cannot reach registry: %s", err.Error())
	}
	defer func() { _ = k.Close() }()

	ki, err := k.Stat()
	if err != nil {
		return nil, fmt.Errorf("cannot stat registry key: %s", err.Error())
	}

	s, err := k.ReadValueNames(int(ki.ValueCount))
	if err != nil {
		return nil, fmt.Errorf("cannot reach registry values: %s", err.Error())
	}

	ret := []string{}
	for _, elt := range s {
		q, _, err := k.GetStringValue(elt)
		if err != nil {
			return nil, fmt.Errorf("cannot extract string value: %s", err.Error())
		}

		if pattern != nil {
			if pattern.MatchString(elt) {
				log.Infof("serial USB detected on %s", q)
				ret = append(ret, q)
			}
			continue
		}

		ret = append(ret, q)
	}

	return ret, nil
}
