package oneandone

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

var (
	set_baremetal_server  sync.Once
	baremetal_server_id   string
	baremetal_server_name string
	baremetal_ser_app_id  string
	baremetal_server      *Server
	baremetal_ser_lb      *LoadBalancer
	baremetalModel        BaremetalModel
)

func setup_baremetal_server() {
	fmt.Println("Deploying a baremetal test server...")
	b_srv_id, b_srv, err := create_baremetal_test_server(false)
	//b_srv, err := api.GetServer("86884513A78F2F2223E393BE1CA1135A")
	//b_srv_id := b_srv.Id

	if err != nil {
		fmt.Printf("Unable to create the baremetal server. Error: %s", err.Error())
		return
	}
	if b_srv_id == "" || b_srv.Id == "" {
		fmt.Printf("Unable to create the baremetal server.")
		return
	} else {
		baremetal_server_id = b_srv.Id
	}

	err = api.WaitForState(b_srv, "POWERED_OFF", 20, 90)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	baremetal_server = b_srv
}

func get_baremetal_appliance(max_disk_size int) ServerAppliance {
	saps, _ := api.ListServerAppliances(1, 500, "", "baremetal", "")
	return saps[0]
}

func create_baremetal_test_server(power_on bool) (string, *Server, error) {
	rand.Seed(time.Now().UnixNano())
	baremetal_server_name = fmt.Sprintf("GOBAREMETAL_%d", rand.Intn(1000000))
	fmt.Printf("Creating test server '%s'...\n", baremetal_server_name)

	sap := get_baremetal_appliance(hdd_size)
	baremetal_ser_app_id = sap.Id
	//mp := get_default_mon_policy()
	baremetalModels, err := api.ListBaremetalModels(1, 1, "", "BMC_L", "")
	baremetalModel = baremetalModels[0]
	baremetalModelId := baremetalModels[0].Id

	req := ServerRequest{
		Name:        baremetal_server_name,
		Description: baremetal_server_name + " description",
		ApplianceId: baremetal_ser_app_id,
		//MonitoringPolicyId: mp.Id,
		PowerOn:    power_on,
		ServerType: "baremetal",
		Hardware: Hardware{
			BaremetalModelId: baremetalModelId,
		},
	}
	b_ser_id, b_server, err := api.CreateServer(&req)
	return b_ser_id, b_server, err
}

// /servers tests

func TestCreateBaremetalServer(t *testing.T) {
	set_baremetal_server.Do(setup_baremetal_server)

	if baremetal_server == nil {
		t.Errorf("CreateServer failed.")
		return
	}
	if baremetal_server.Name != baremetal_server_name {
		t.Errorf("Wrong server name.")
	}
	if baremetal_server.Image.Id != ser_app_id {
		t.Errorf("Wrong server image on server '%s'.", baremetal_server.Name)
	}
}

func TestListAllServers(t *testing.T) {
	//set_baremetal_server.Do(setup_baremetal_server)
	fmt.Println("Listing all servers...")

	res, err := api.ListServers()
	if err != nil {
		t.Errorf("ListServers failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No server found.")
	}

	res, err = api.ListServers(1, 2, "name", "", "id,name")

	if err != nil {
		t.Errorf("ListServers with parameter options failed. Error: " + err.Error())
		return
	}
	if len(res) == 0 {
		t.Errorf("No server found.")
	}
	if len(res) > 2 {
		t.Errorf("Wrong number of objects per page.")
	}
	if res[0].Hardware != nil {
		t.Errorf("Filtering parameters failed.")
	}
	if res[0].Name == "" {
		t.Errorf("Filtering parameters failed.")
	}
	if len(res) == 2 && res[0].Name >= res[1].Name {
		t.Errorf("Sorting parameters failed.")
	}
	// Test for error response
	res, err = api.ListServers(0, 0, true, "name", "")
	if res != nil || err == nil {
		t.Errorf("ListServers failed to handle incorrect argument type.")
	}

	res, err = api.ListServers(0, 0, "", baremetal_server_name, "")

	if err != nil {
		t.Errorf("ListServers with parameter options failed. Error: " + err.Error())
		return
	}
	if len(res) != 1 {
		t.Errorf("Search parameter failed.")
	}
	if res[0].Name != baremetal_server_name {
		t.Errorf("Search parameter failed.")
	}
}

func TestGetBaremetalServer(t *testing.T) {
	set_baremetal_server.Do(setup_baremetal_server)

	fmt.Println("Getting the server...")
	b_srv, err := api.GetServer(baremetal_server_id)

	if err != nil {
		t.Errorf("GetServer failed. Error: " + err.Error())
		return
	}
	if b_srv.Id != baremetal_server_id {
		t.Errorf("Wrong server ID.")
	}
}

func TestGetBaremetalServerStatus(t *testing.T) {
	set_baremetal_server.Do(setup_baremetal_server)

	fmt.Println("Getting the server's status...")
	status, err := api.GetServerStatus(baremetal_server_id)

	if err != nil {
		t.Errorf("GetServerStatus failed. Error: " + err.Error())
		return
	}
	if status.State != "POWERED_OFF" {
		t.Errorf("Wrong server status.")
	}
}

func TestStartBaremetalServer(t *testing.T) {
	set_baremetal_server.Do(setup_baremetal_server)

	fmt.Println("Starting the server...")
	b_srv, err := api.StartServer(baremetal_server_id)

	if err != nil {
		t.Errorf("StartServer failed. Error: " + err.Error())
		return
	}

	err = api.WaitForState(b_srv, "POWERED_ON", 10, 60)

	if err != nil {
		t.Errorf("Starting the server failed. Error: " + err.Error())
	}

	server, err = api.GetServer(b_srv.Id)

	if err != nil {
		t.Errorf("GetServer failed. Error: " + err.Error())
	} else if server.Status.State != "POWERED_ON" {
		t.Errorf("Wrong server state. Expected: POWERED_ON. Found: %s.", server.Status.State)
	}
}

func TestRebootBaremetalServer(t *testing.T) {
	set_baremetal_server.Do(setup_baremetal_server)

	for i := 1; i < 3; i++ {
		is_hardware := i%2 == 0
		var method string
		if is_hardware {
			method = "HARDWARE"
		} else {
			method = "SOFTWARE"
		}
		fmt.Printf("Rebooting the server using '%s' method...\n", method)
		b_srv, err := api.RebootServer(baremetal_server_id, is_hardware)

		if err != nil {
			t.Errorf("RebootServer using '%s' method failed. Error: %s", method, err.Error())
			return
		}

		err = api.WaitForState(b_srv, "REBOOTING", 10, 60)

		if err != nil {
			t.Errorf("Rebooting the server using '%s' method failed. Error:  %s", method, err.Error())
		}

		err = api.WaitForState(b_srv, "POWERED_ON", 10, 60)

		if err != nil {
			t.Errorf("Rebooting the server using '%s' method failed. Error:  %s", method, err.Error())
		}
	}
}

func TestRenameBaremetalServer(t *testing.T) {
	set_baremetal_server.Do(setup_baremetal_server)

	fmt.Println("Renaming the server...")

	new_name := server.Name + "_renamed"
	new_desc := server.Description + "_renamed"

	b_srv, err := api.RenameServer(baremetal_server_id, new_name, new_desc)

	if err != nil {
		t.Errorf("Renaming server failed. Error: " + err.Error())
		return
	}
	if b_srv.Name != new_name {
		t.Errorf("Wrong server name.")
	}
	if b_srv.Description != new_desc {
		t.Errorf("Wrong server description.")
	}
}

func TestGetBaremetalServerHardware(t *testing.T) {
	set_baremetal_server.Do(setup_baremetal_server)

	fmt.Println("Getting the server's hardware...")
	hardware, err := api.GetServerHardware(baremetal_server_id)
	baremetalModel, err := api.GetBaremetalModel("81504C620D98BCEBAA5202D145203B4B")

	if err != nil {
		t.Errorf("GetServerHardware failed. Error: " + err.Error())
		return
	}
	if hardware == nil {
		t.Errorf("Unable to get the server's hardware.")
	} else {
		if hardware.Vcores != baremetalModel.Hardware.CoresPerProcessor {
			t.Errorf("Wrong number of processor cores on server '%s'.", server.Name)
		}
		if hardware.CoresPerProcessor != baremetalModel.Hardware.Cores {
			t.Errorf("Wrong number of cores per processor on server '%s'.", server.Name)
		}
		if hardware.Ram != baremetalModel.Hardware.Ram {
			t.Errorf("Wrong RAM size on server '%s'.", server.Name)
		}
	}
}

func TestListBaremetalServerHdds(t *testing.T) {
	set_baremetal_server.Do(setup_baremetal_server)

	fmt.Println("Listing all the server's HDDs...")
	hdds, err := api.ListServerHdds(baremetal_server_id)

	if err != nil {
		t.Errorf("ListServerHdds failed. Error: " + err.Error())
		return
	}
	if len(hdds) != 1 {
		t.Errorf("Wrong number of the server's hard disks.")
	}
	if hdds[0].Id == "" {
		t.Errorf("Wrong HDD id.")
	}
	if hdds[0].Size != 1600 {
		t.Errorf("Wrong HDD size.")
	}
	if hdds[0].IsMain != true {
		t.Errorf("Wrong main HDD.")
	}
	server_hdd = &hdds[0]
}

func TestGetBaremetalServerHdd(t *testing.T) {
	set_baremetal_server.Do(setup_baremetal_server)
	hdds, _ := api.ListServerHdds(baremetal_server_id)

	fmt.Println("Getting server HDD...")
	hdd, err := api.GetServerHdd(baremetal_server_id, hdds[0].Id)

	if err != nil {
		t.Errorf("GetServerHdd failed. Error: " + err.Error())
		return
	}
	if hdd.Id != hdds[0].Id {
		t.Errorf("Wrong HDD id.")
	}
	if hdd.Size != hdds[0].Size {
		t.Errorf("Wrong HDD size.")
	}
	if hdd.IsMain != hdds[0].IsMain {
		t.Errorf("Wrong main HDD.")
	}
}

func TestGetBaremetalServerImage(t *testing.T) {
	set_baremetal_server.Do(setup_baremetal_server)

	fmt.Println("Getting the server's image...")
	img, err := api.GetServerImage(baremetal_server_id)

	if err != nil {
		t.Errorf("GetServerImage failed. Error: " + err.Error())
		return
	}
	if img.Id != baremetal_server.Image.Id {
		t.Errorf("Wrong image ID.")
	}
	if img.Name != baremetal_server.Image.Name {
		t.Errorf("Wrong image name.")
	}
}

func TestStopBaremetalServer(t *testing.T) {
	set_baremetal_server.Do(setup_baremetal_server)

	fmt.Println("Stopping the server...")
	b_srv, err := api.ShutdownServer(baremetal_server_id, false)

	if err != nil {
		t.Errorf("ShutdownServer failed. Error: " + err.Error())
	} else {
		err = api.WaitForState(b_srv, "POWERED_OFF", 10, 60)
		if err != nil {
			t.Errorf("Stopping the server failed. Error: " + err.Error())
		}

		baremetal_server, err = api.GetServer(baremetal_server_id)
		if err != nil {
			t.Errorf("GetServer failed. Error: " + err.Error())
		}
		if baremetal_server.Status.State != "POWERED_OFF" {
			t.Errorf("Wrong baremetal_server state. Expected: POWERED_OFF. Found: %s.", baremetal_server.Status.State)
		}
	}
}

func TestListBaremetalServerIps(t *testing.T) {
	set_baremetal_server.Do(setup_baremetal_server)

	b_srv, e := api.GetServer(baremetal_server_id)
	if e == nil {
		baremetal_server = b_srv
	}

	fmt.Println("Listing the baremetal_server's IPs...")
	ips, err := api.ListServerIps(baremetal_server_id)

	if err != nil {
		t.Errorf("ListServerIps failed. Error: " + err.Error())
		return
	}
	if len(ips) != len(baremetal_server.Ips) {
		t.Errorf("Not all IPs were obtained.")
	}
	if ips[0].Id != baremetal_server.Ips[0].Id {
		t.Errorf("Wrong IP ID.")
	}
	if ips[0].Ip != baremetal_server.Ips[0].Ip {
		t.Errorf("Wrong IP address.")
	}
}

func TestAssignBaremetalServerIp(t *testing.T) {
	set_baremetal_server.Do(setup_baremetal_server)

	fmt.Println("Assigning new IP addresses to the baremetal_server...")
	for i := 2; i < 4; i++ {
		time.Sleep(time.Second)
		b_srv, err := api.AssignServerIp(baremetal_server_id, "IPV4")
		if err != nil {
			t.Errorf("AssignServerIp failed. Error: " + err.Error())
			return
		}
		b_srv = wait_for_action_done(b_srv, 10, 30)
		if len(b_srv.Ips) != i {
			t.Errorf("IP address not assigned to the baremetal_server.")
		}
		baremetal_server = b_srv
	}
}

func TestGetBaremetalServerIp(t *testing.T) {
	set_baremetal_server.Do(setup_baremetal_server)

	fmt.Println("Getting the baremetal_server's IP...")
	baremetal_server, _ = api.GetServer(baremetal_server_id)
	if baremetal_server == nil {
		t.Errorf("GetServer failed.")
		return
	}
	time.Sleep(time.Second)
	ip, err := api.GetServerIp(baremetal_server_id, baremetal_server.Ips[0].Id)

	if err != nil {
		t.Errorf("GetServerIps failed. Error: " + err.Error())
		return
	}
	if ip.Id != baremetal_server.Ips[0].Id {
		t.Errorf("Wrong IP ID.")
	}
	if ip.Ip != baremetal_server.Ips[0].Ip {
		t.Errorf("Wrong IP address.")
	}
}

func TestDeleteBaremetalServerIp(t *testing.T) {
	set_baremetal_server.Do(setup_baremetal_server)

	if len(baremetal_server.Ips) <= 1 {
		for i := 0; i < 2; i++ {
			time.Sleep(10 * time.Second)
			s, e := api.AssignServerIp(baremetal_server_id, "IPV4")
			if s != nil && e == nil {
				s = wait_for_action_done(s, 10, 30)
				baremetal_server = s
				time.Sleep(120 * time.Second)
			}
		}
	}
	ip_no := len(baremetal_server.Ips)
	for i := 1; i < ip_no-1; i++ {
		keep_ip := i%2 == 0
		fmt.Printf("Deleting the baremetal_server's IP '%s' (keep_ip = %s)...\n", baremetal_server.Ips[i].Ip, strconv.FormatBool(keep_ip))
		api.DeleteServerIp(baremetal_server_id, baremetal_server.Ips[i].Id, keep_ip)
		time.Sleep(280 * time.Second)
		ip, _ := api.GetPublicIp(baremetal_server.Ips[i].Id)
		if keep_ip {
			if ip == nil {
				t.Errorf("Failed to keep public IP '%s' when removed from baremetal_server.", baremetal_server.Ips[i].Ip)
			} else {
				fmt.Printf("Deleting IP address '%s' after removing from the baremetal_server...\n", baremetal_server.Ips[i].Ip)
				api.DeletePublicIp(ip.Id)
			}
		} else if ip != nil {
			t.Errorf("Failed to delete public IP '%s' when removed from baremetal_server.", baremetal_server.Ips[i].Ip)
			fmt.Printf("Cleaning up. Deleting IP address '%s' directly...\n", baremetal_server.Ips[i].Ip)
			api.DeletePublicIp(ip.Id)
		}
	}
}

func TestAssignBaremetalServerIpLoadBalancer(t *testing.T) {
	set_baremetal_server.Do(setup_baremetal_server)
	ips, _ := api.ListServerIps(baremetal_server_id)
	lb := create_load_balancer()

	fmt.Printf("Assigning a load balancer to the baremetal_server's IP '%s'...\n", ips[0].Ip)
	b_srv, err := api.AssignServerIpLoadBalancer(baremetal_server_id, ips[0].Id, lb.Id)

	if err != nil {
		t.Errorf("AssignServerIpLoadBalancer failed. Error: " + err.Error())
		return
	}
	if len(b_srv.Ips[0].LoadBalancers) == 0 {
		t.Errorf("Load balancer not assigned.")
	}
	if b_srv.Ips[0].LoadBalancers[0].Id != lb.Id {
		t.Errorf("Wrong load balancer assigned.")
	}
	baremetal_ser_lb = lb
}

func TestListBaremetalServerIpLoadBalancers(t *testing.T) {
	set_baremetal_server.Do(setup_baremetal_server)
	ips, _ := api.ListServerIps(baremetal_server_id)

	fmt.Println("Listing load balancers assigned to the baremetal_server's IP...")
	lbs, err := api.ListServerIpLoadBalancers(baremetal_server_id, ips[0].Id)

	if err != nil {
		t.Errorf("ListServerIpLoadBalancers failed. Error: " + err.Error())
		return
	}
	if len(lbs) == 0 {
		t.Errorf("No load balancer was assigned to the baremetal_server's IP.")
		return
	}
	if lbs[0].Id != baremetal_ser_lb.Id {
		t.Errorf("Wrong load balancer assigned.")
	}
}

func TestUnassignBaremetalServerIpLoadBalancer(t *testing.T) {
	set_baremetal_server.Do(setup_baremetal_server)
	ips, _ := api.ListServerIps(baremetal_server_id)

	fmt.Println("Unassigning the load balancer from the baremetal_server's IP...")
	b_srv, err := api.UnassignServerIpLoadBalancer(baremetal_server_id, ips[0].Id, baremetal_ser_lb.Id)

	if err != nil {
		t.Errorf("UnassignServerIpLoadBalancer failed. Error: " + err.Error())
		return
	}
	if len(b_srv.Ips[0].LoadBalancers) > 0 {
		t.Errorf("Unassigning the load balancer failed.")
	}
	baremetal_ser_lb, err = api.DeleteLoadBalancer(baremetal_ser_lb.Id)
	if err == nil {
		api.WaitUntilDeleted(baremetal_ser_lb)
	}
	baremetal_ser_lb, _ = api.GetLoadBalancer(baremetal_ser_lb.Id)
}

func TestAssignBaremetalServerIpFirewallPolicy(t *testing.T) {
	set_baremetal_server.Do(setup_baremetal_server)
	ips, _ := api.ListServerIps(baremetal_server_id)

	fmt.Println("Assigning a firewall policy to the baremetal_server's IP...")
	fps, err := api.ListFirewallPolicies(0, 1, "creation_date", "linux", "id,name")
	if err != nil {
		t.Errorf("ListFirewallPolicies failed. Error: " + err.Error())
		return
	}
	b_srv, err := api.AssignServerIpFirewallPolicy(baremetal_server_id, ips[0].Id, fps[0].Id)

	if err != nil {
		t.Errorf("AssignServerIpFirewallPolicy failed. Error: " + err.Error())
		return
	}
	if b_srv.Ips[0].Firewall == nil {
		t.Errorf("Firewall policy not assigned.")
	}
	if b_srv.Ips[0].Firewall.Id != fps[0].Id {
		t.Errorf("Wrong firewall policy assigned.")
	}
}

func TestGetBaremetalServerIpFirewallPolicy(t *testing.T) {
	set_baremetal_server.Do(setup_baremetal_server)
	ips, _ := api.ListServerIps(baremetal_server_id)

	fmt.Println("Getting the firewall policy assigned to the baremetal_server's IP...")
	fps, err := api.ListFirewallPolicies(0, 1, "creation_date", "linux", "id,name")
	if err != nil {
		t.Errorf("ListFirewallPolicies failed. Error: " + err.Error())
	}
	fp, err := api.GetServerIpFirewallPolicy(baremetal_server_id, ips[0].Id)

	if err != nil {
		t.Errorf("GetServerIpFirewallPolicy failed. Error: " + err.Error())
	}
	if fp == nil {
		t.Errorf("No firewall policy assigned to the baremetal_server's IP.")
	}
	if fp.Id != fps[0].Id {
		t.Errorf("Wrong firewall policy assigned to the baremetal_server's IP.")
	}
}

func TestDeleteBaremetalServer(t *testing.T) {
	set_baremetal_server.Do(setup_baremetal_server)

	time.Sleep(120 * time.Second)
	b_srv, err := api.DeleteServer(baremetal_server_id, true)
	if err != nil {
		t.Errorf("DeleteServer server failed. Error: " + err.Error())
		return
	}
	fmt.Printf("Deleting baremetal_server '%s', keeping baremetal_server's IP '%s'...\n", b_srv.Name, b_srv.Ips[0].Ip)
	ip_id := b_srv.Ips[0].Id

	err = api.WaitUntilDeleted(b_srv)

	if err != nil {
		t.Errorf("Deleting the baremetal_server failed. Error: " + err.Error())
	}

	b_srv, err = api.GetServer(baremetal_server_id)

	if b_srv != nil {
		t.Errorf("Unable to delete the baremetal_server.")
	} else {
		baremetal_server = nil
	}

	ip, _ := api.GetPublicIp(ip_id)
	if ip == nil {
		t.Errorf("Failed to keep IP after deleting the baremetal_server.")
	} else {
		fmt.Printf("Deleting baremetal_server's IP '%s' after deleting the baremetal_server...\n", ip.IpAddress)
		ip, err = api.DeletePublicIp(ip_id)
		if err != nil {
			t.Errorf("Unable to delete baremetal_server's IP after deleting the baremetal_server.")
		} else {
			api.WaitUntilDeleted(ip)
		}
	}
}
