package udisks

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus"
)

func getPropertyString(conn *dbus.Conn, path dbus.ObjectPath, propertyName string) (string, error) {
	obj := conn.Object("org.freedesktop.UDisks2", path)
	result, err := obj.GetProperty(propertyName)
	if err != nil {
		return "", err
	}
	return result.Value().(string), nil
}

func getPropertyStringArray(conn *dbus.Conn, path dbus.ObjectPath, propertyName string) ([]string, error) {
	obj := conn.Object("org.freedesktop.UDisks2", path)
	p, err := obj.GetProperty(propertyName)
	if err != nil {
		return nil, err
	}
	if byteArray, ok := p.Value().([][]byte); ok {
		result := make([]string, 0)
		for _, bytes := range byteArray {
			result = append(result, strings.TrimRight(string(bytes), "\x00"))
		}
		return result, nil
	}
	return nil, fmt.Errorf("invalid property type")
}
