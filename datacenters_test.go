package oneandone

import (
	"fmt"
	"testing"
	"time"
)

// /datacenters tests

func TestListDatacenters(t *testing.T) {
	fmt.Println("Listing datacenters...")

	res, err := api.ListDatacenters()
	if err != nil {
		t.Errorf("ListDatacenters failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No datacenter found.")
	}

	res, err = api.ListDatacenters(1, 2, "location", "", "id,location")

	if err != nil {
		t.Errorf("ListDatacenters with parameter options failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No datacenter found.")
	}

	// Test for error response
	res, err = api.ListDatacenters(0, 0, "location", "Spain", "id", "country_code")
	if res != nil || err == nil {
		t.Errorf("ListDatacenters failed to handle incorrect number of passed arguments.")
	}

	res, err = api.ListDatacenters(0, 0, "", "Germany", "")

	if err != nil {
		t.Errorf("ListDatacenters with parameter options failed. Error: " + err.Error())
	}
}

func TestGetDatacenter(t *testing.T) {
	dcs, err := api.ListDatacenters()

	if len(dcs) == 0 {
		t.Errorf("No datacenter found. " + err.Error())
		return
	}

	for i, _ := range dcs {
		time.Sleep(time.Second)
		fmt.Printf("Getting datacenter '%s'...\n", dcs[i].CountryCode)
		dc, err := api.GetDatacenter(dcs[i].Id)

		if err != nil {
			t.Errorf("GetDatacenter failed. Error: " + err.Error())
			return
		}
		if dc.Id != dcs[i].Id {
			t.Errorf("Wrong datacenter ID.")
		}
		if dc.CountryCode != dcs[i].CountryCode {
			t.Errorf("Wrong country code of the datacenter.")
		}
		if dc.Location != dcs[i].Location {
			t.Errorf("Wrong datacenter location.")
		}
	}
}
