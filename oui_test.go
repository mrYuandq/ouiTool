package OUI

import (
	"testing"
	"fmt"
	"time"
	"os"
)

var Oui *OuiDb
func TestMain(m *testing.M) {
	setup()
	time1 := time.Now()
	code := m.Run()
	fmt.Println("use timed:", time.Now().Sub(time1))
	shutDown()
	os.Exit(code)
}
func setup() {
	Oui = &OuiDb{}
	fmt.Println("setup tests")
}
func shutDown() {
	fmt.Println("shutDown tests")
}

func TestOuiDb_SetInsideOrganization(t *testing.T) {
	// must be Capitalized
	Oui.SetInsideOrganization([]string{
		"APPLE",
	})
	Oui = Oui.Open("../oui.txt")
	if nil == Oui {
		fmt.Println("open oui.txt failed")
		return
	}

	// only APPLE Organization can find
	addr, err := Oui.VendorLookup("00254B")
	if nil != err {
		fmt.Println("err:", err.Error())
		return
	}
	fmt.Println(addr.Organization)

	addr, err = Oui.VendorLookup("C8D779")
	if nil != err {
		fmt.Println("err:", err.Error())
		return
	}
	fmt.Println(addr.Organization)
}

func TestOuiDb_SetOutsideOrganization(t *testing.T) {
	// must be Capitalized
	Oui.SetOutsideOrganization([]string{
		"APPLE",
	})
	Oui = Oui.Open("../oui.txt")
	if nil == Oui {
		fmt.Println("open oui.txt failed")
		return
	}

	addr, err := Oui.VendorLookup("C8D779")
	if nil != err {
		fmt.Println("err:", err.Error())
		return
	}
	fmt.Println(addr.Organization)

	addr, err = Oui.VendorLookup("00254B")
	if nil != err {
		fmt.Println("err:", err.Error())
		return
	}
	fmt.Println(addr.Organization)
}

func TestOuiDb_VendorLookup(t *testing.T) {
	Oui = Oui.Open("../oui.txt")
	if nil == Oui {
		fmt.Println("open oui.txt failed")
		return
	}

	addr, err := Oui.VendorLookup("00254B")
	if nil != err {
		fmt.Println("err:", err.Error())
		return
	}
	fmt.Println(addr.Organization)
}

