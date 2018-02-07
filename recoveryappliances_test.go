package oneandone

import (
	"fmt"
	"strings"
	"testing"
)

// /recovery_appliances tests

func TestListRecoveryAppliances(t *testing.T) {
	fmt.Println("Listing all recovery appliances...")

	res, err := api.ListRecoveryAppliances()
	if err != nil {
		t.Errorf("ListRecoveryAppliances failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No server appliance found.")
	}

	res, err = api.ListRecoveryAppliances(1, 2, "name", "", "id,name")

	if err != nil {
		t.Errorf("ListRecoveryAppliances with parameter options failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No recovery appliance found.")
	}
	if len(res) != 2 {
		t.Errorf("Wrong number of objects per page.")
	}
	for index := 0; index < len(res); index += 1 {
		if res[index].Id == "" {
			t.Errorf("Filtering a list of recovery appliances failed.")
		}
		if res[index].Name == "" {
			t.Errorf("Filtering a list of recovery appliances failed.")
		}

		if index < len(res)-1 {
			if res[index].Name > res[index+1].Name {
				t.Errorf("Sorting a list of recovery appliances failed.")
			}
		}
	}
	// Test for error response
	res, err = api.ListRecoveryAppliances(nil, nil, nil, nil, nil)
	if res != nil || err == nil {
		t.Errorf("ListRecoveryAppliances failed to handle incorrect argument type.")
	}

	res, err = api.ListRecoveryAppliances(0, 0, "", "linux", "")

	if err != nil {
		t.Errorf("ListRecoveryAppliances with parameter options failed. Error: " + err.Error())
	}

	for _, sa := range res {
		if !strings.Contains(strings.ToLower(sa.Os.Name), "linux") {
			t.Errorf("Search parameter failed.")
		}
	}
}

func TestGetRecoveryAppliance(t *testing.T) {
	raps, _ := api.ListRecoveryAppliances(1, 1, "", "", "")
	fmt.Printf("Getting recovery appliance '%s'...\n", raps[0].Name)
	sa, err := api.GetRecoveryAppliance(raps[0].Id)

	if sa == nil || err != nil {
		t.Errorf("GetRecoveryAppliance failed. Error: " + err.Error())
	}
	if sa.Id != raps[0].Id {
		t.Errorf("Wrong ID of the recovery appliance.")
	}
}
