package collector

import "testing"

func TestHello(t *testing.T) {
	if Hello() == "Hello" {
		t.Log("collector.Hello PASS")
	} else {
		t.Error("collector.Hello Failed")
	}
}
