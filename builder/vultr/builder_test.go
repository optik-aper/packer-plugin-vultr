package vultr

import (
	"strconv"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/packer"
)

func testConfig() map[string]interface{} {
	return map[string]interface{}{
		"api_key":              "test-api-key",
		"snapshot_description": "packer-test-snapshot",
		"region_id":            "ewr",
		"os_id":                "352",
		"plan_id":              "vc2-1c-1gb",
		"ssh_username":         "root",
	}
}

func TestBuilder_ImplementsBuilder(t *testing.T) {
	var raw interface{} = &Builder{}
	if _, ok := raw.(packer.Builder); !ok {
		t.Fatalf("Builder should be a builder")
	}
}

func TestBuilder_Prepare_BadType(t *testing.T) {
	b := &Builder{}
	c := map[string]interface{}{
		"api_key": []string{},
	}

	_, warnings, err := b.Prepare(c)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err == nil {
		t.Fatalf("prepare should fail")
	}
}

func TestBuilderPrepare_InvalidKey(t *testing.T) {
	var b Builder
	config := testConfig()

	// Add a random key
	config["i_should_not_be_valid"] = true
	_, warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err == nil {
		t.Fatal("should have error")
	}
}

func TestBuilderPrepare_RegionID(t *testing.T) {
	var b Builder
	config := testConfig()

	// Test default
	delete(config, "region_id")
	_, warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err == nil {
		t.Fatalf("should error")
	}

	expected := "ewr"

	// Test set
	config["region_id"] = expected
	b = Builder{}
	_, warnings, err = b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	if b.config.RegionID != expected {
		t.Errorf("found %s, expected %s", b.config.RegionID, expected)
	}
}

func TestBuilderPrepare_OSID(t *testing.T) {
	var b Builder
	config := testConfig()

	// Test default
	delete(config, "os_id")
	_, warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should error")
	}

	expected := 352

	// Test set
	config["os_id"] = expected
	b = Builder{}
	_, warnings, err = b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	if b.config.OSID != expected {
		t.Errorf("found %d, expected %d", b.config.OSID, expected)
	}
}

func TestBuilderPrepare_PlanID(t *testing.T) {
	var b Builder
	config := testConfig()

	// Test default
	delete(config, "plan_id")
	_, warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err == nil {
		t.Fatalf("should error")
	}

	expected := "vc2-1c-1gb"

	// Test set
	config["plan_id"] = expected
	b = Builder{}
	_, warnings, err = b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	if b.config.PlanID != expected {
		t.Errorf("found %s, expected %s", b.config.PlanID, expected)
	}
}

func TestBuilderPrepare_Description(t *testing.T) {
	var b Builder
	config := testConfig()

	// Test default
	_, warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	if b.config.Description == "" {
		t.Errorf("invalid: %s", b.config.Description)
	}

	// Test set
	config["snapshot_description"] = "foobarbaz"
	b = Builder{}
	_, warnings, err = b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	// Test set with template
	config["snapshot_description"] = "{{timestamp}}"
	b = Builder{}
	_, warnings, err = b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	_, err = strconv.ParseInt(b.config.Description, 0, 0)
	if err != nil {
		t.Fatalf("failed to parse int in template: %s", err)
	}
}

func TestBuilderPrepare_StateTimeout(t *testing.T) {
	var b Builder
	config := testConfig()

	// Test default
	_, warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	// Test set
	config["state_timeout"] = "5m"
	b = Builder{}
	_, warnings, err = b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	// Test bad
	config["state_timeout"] = "tubes"
	b = Builder{}
	_, warnings, err = b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err == nil {
		t.Fatal("should have error")
	}
}

func TestBuilderPrepare_EnablePrivateNetwork(t *testing.T) {
	var b Builder
	config := testConfig()

	// Test default
	_, warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	if b.config.EnablePrivateNetwork != false {
		t.Errorf("invalid: %t", b.config.EnablePrivateNetwork)
	}

	// Test set
	config["enable_private_network"] = true
	b = Builder{}
	_, warnings, err = b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	if b.config.EnablePrivateNetwork != true {
		t.Errorf("invalid: %t", b.config.EnablePrivateNetwork)
	}
}

func TestBuilderPrepare_EnableIPV6(t *testing.T) {
	var b Builder
	config := testConfig()

	// Test default
	_, warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	if b.config.EnableIPV6 != false {
		t.Errorf("invalid: %t", b.config.EnableIPV6)
	}

	// Test set
	config["enable_ipv6"] = true
	b = Builder{}
	_, warnings, err = b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	if b.config.EnableIPV6 != true {
		t.Errorf("invalid: %t", b.config.EnableIPV6)
	}
}

func TestBuilderPrepare_Label(t *testing.T) {
	var b Builder
	config := testConfig()

	// Test default
	_, warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	if b.config.Label == "" {
		t.Errorf("invalid: %s", b.config.Label)
	}

	// Test set
	config["instance_label"] = "foobarbaz"
	b = Builder{}
	_, warnings, err = b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	// Test set with template
	config["instance_label"] = "{{timestamp}}"
	b = Builder{}
	_, warnings, err = b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	_, err = strconv.ParseInt(b.config.Label, 0, 0)
	if err != nil {
		t.Fatalf("failed to parse int in template: %s", err)
	}
}
