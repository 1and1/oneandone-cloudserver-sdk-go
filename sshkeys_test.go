package oneandone

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	set_ssh_key   sync.Once
	test_ssh_key  *SSHKey
	test_ssh_name string
	test_ssh_desc string
)

func setup_ssh_key() {
	rand.Seed(time.Now().UnixNano())
	rint := rand.Intn(999)
	fmt.Printf("Creating test ssh_key '%s'...\n", fmt.Sprintf("SSHKEY_%d", rint))

	pk := `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCezYMOgAan+JmatgFJ+Q1FUNjrqNVgAvcTkjYJwHVcQaolq9f9qB7tEeUPDNj2oNN49joAmTcllDcPIxryT5PnQUaaUhu4ZJ9+bRtXCyhnf2LJQdVfzBFEBJX9fW4RiV1XtSAtLRBrbrCb4JjHmhIYpvhBHC29Ve+g64nhdvBhqyLZ3SLI2U/opEmt5u2xftWGl0TBSQYveqc4ntz3fe+f9XlBHvK3Nw12bCLmLle7jQuZ4lXyAYqNAfdOMTs2zMTk422Dl/h4+zRh1h4rM9zaCk4+g3kdugJm7Vul03wm43cHmHsJv51R3XKSHzgb7q/eNj+YdMi5Ndt0Bm+bLjw7`

	req := SSHKeyRequest{
		Name:        fmt.Sprintf("SSHKEY_%d", rint),
		Description: fmt.Sprintf("SSHKEY_%d description", rint),
		PublicKey:   pk,
	}

	sshkey_id, sshkey, err := api.CreateSSHKey(&req)

	if err != nil {
		fmt.Printf("Unable to create the ssh key. Error: %s", err.Error())
		return
	}
	if sshkey_id == "" || sshkey.Id == "" {
		fmt.Printf("Unable to create the server.")
		return
	} else {
		sshkey_id = sshkey.Id
	}

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	test_ssh_key = sshkey
}

func TestListSSHKeys(t *testing.T) {
	set_ssh_key.Do(setup_ssh_key)
	keys, err := api.ListSSHKeys()

	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}

	if len(keys) < 1 {
		t.Errorf("Assertion failed.")
	}
}

func TestCreateSSHKey(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	rint := rand.Intn(999)
	pk := `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCezYMOgAan+JmatgFJ+Q1FUNjrqNVgAvcTkjYJwHVcQaolq9f9qB7tEeUPDNj2oNN49joAmTcllDcPIxryT5PnQUaaUhu4ZJ9+bRtXCyhnf2LJQdVfzBFEBJX9fW4RiV1XtSAtLRBrbrCb4JjHmhIYpvhBHC29Ve+g64nhdvBhqyLZ3SLI2U/opEmt5u2xftWGl0TBSQYveqc4ntz3fe+f9XlBHvK3Nw12bCLmLle7jQuZ4lXyAYqNAfdOMTs2zMTk422Dl/h4+zRh1h4rM9zaCk4+g3kdugJm7Vul03wm43cHmHsJv51R3XKSHzgb7q/eNj+YdMi5Ndt0Bm+bLjw7`

	req := SSHKeyRequest{
		Name:        fmt.Sprintf("SSHKEY_%d", rint),
		Description: fmt.Sprintf("SSHKEY_%d description", rint),
		PublicKey:   pk,
	}

	_, sshkey, err := api.CreateSSHKey(&req)

	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}

	if sshkey.PublicKey != strings.Replace(pk, "ssh-rsa ", "", -1) {
		t.Errorf("Keys are not the same")
		t.Fail()
	}

	_, err = api.DeleteSSHKey(sshkey.Id)
	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}
}

func TestRenameSSHKey(t *testing.T) {
	set_ssh_key.Do(setup_ssh_key)

	sshKey, err := api.RenameSSHKey(test_ssh_key.Id, test_ssh_key.Name+"1", test_ssh_key.Description+"1")

	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}

	if sshKey.Name == sshKey.Name+"1" {
		t.Errorf("Names are not the same")
		t.Fail()
	}

	if sshKey.Description == sshKey.Description+"1" {
		t.Errorf("Descriptions are not the same")
		t.Fail()
	}
}

func TestDeleteSSHKey(t *testing.T) {
	set_ssh_key.Do(setup_ssh_key)
	_, err := api.DeleteSSHKey(test_ssh_key.Id)
	if err != nil {
		t.Errorf("Error while deleting ssh key", err.Error())
		t.Fail()
	}

	sshKey, err := api.GetSSHKey(test_ssh_key.Id)

	if err == nil {
		t.Errorf(err.Error())
		t.Fail()
	}

	if sshKey != nil {
		t.Errorf("SSH Key was not deleted")
		t.Fail()
	}

}
