package network

import (
	"testing"
)

var testName = "testbridge"

func TestBridgeCreate(t *testing.T) {
	d := BridgeNetworkDriver{}
	n, err := d.Create("192.168.0.1/24", testName)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("create network :%v", n)
}

func TestBridgeDelete(t *testing.T) {
	d := BridgeNetworkDriver{}
	err := d.Delete(testName)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("delete network :%v", testName)
}

func TestBridgeConnect(t *testing.T) {
	ep := Endpoint{
		ID: "testcontainer",
	}

	n := Network{
		Name: testName,
	}

	d := BridgeNetworkDriver{}
	err := d.Connect(&n, &ep)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBridgeDisconnect(t *testing.T) {
	ep := Endpoint{
		ID: "testcontainer",
	}

	n := Network{
		Name: testName,
	}

	d := BridgeNetworkDriver{}
	err := d.Disconnect(n, &ep)
	if err != nil {
		t.Fatal(err)
	}
}
