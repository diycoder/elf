package net

import "testing"

func TestGetIp(t *testing.T) {
	ip, err := GetIP()
	t.Log(ip, err)
}
