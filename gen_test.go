package main

import (
	"testing"
	"time"

	v1 "github.com/jiayinzhang-mint/protoc-gen-go-govmomi-bridge/gen/proto/v1"
	"github.com/stretchr/testify/assert"
	vmomiTypes "github.com/vmware/govmomi/vim25/types"
)

// TestHackedPtr tests the case where the field is accidentally a pointer in the govmomi
func TestHackedPtr(t *testing.T) {
	vs := vmomiTypes.WaitOptions{
		MaxWaitSeconds:   v1.NewPointer(int32(10)),
		MaxObjectUpdates: 10,
	}

	ns := &v1.WaitOptions{}

	if err := v1.FromGovmomi(vs, ns); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, *vs.MaxWaitSeconds, ns.MaxWaitSeconds)
	assert.Equal(t, vs.MaxObjectUpdates, *ns.MaxObjectUpdates)
}

// TestSlice tests the case where the field is a slice in the govmomi
func TestSlice(t *testing.T) {
	vs := vmomiTypes.NetworkProfile{
		Vswitch: []vmomiTypes.VirtualSwitchProfile{
			{
				ApplyProfile: vmomiTypes.ApplyProfile{
					Enabled: true,
				},
				Key:  "key-1",
				Name: "name-1",
			},
			{
				ApplyProfile: vmomiTypes.ApplyProfile{
					Enabled: true,
				},
				Key:  "key-21",
				Name: "name-21",
			},
		},
	}

	ns := &v1.NetworkProfile{}

	if err := v1.FromGovmomi(vs, ns); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(vs.Vswitch), len(ns.Vswitch))
	assert.Equal(t, vs.Vswitch[0].ApplyProfile.Enabled, ns.Vswitch[0].Enabled)
}

// TestAnonymous tests the case where the field contains an anonymous struct in the govmomi
func TestAnonymous(t *testing.T) {
	vs := vmomiTypes.NetworkProfile{
		ApplyProfile: vmomiTypes.ApplyProfile{
			Enabled: true,
		},
	}

	ns := &v1.NetworkProfile{}

	if err := v1.FromGovmomi(vs, ns); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, vs.ApplyProfile.Enabled, ns.Enabled)
}

// TestManagedObjectReference tests the case where the field is a ManagedObjectReference in the govmomi, which has weird tag name
func TestManagedObjectReference(t *testing.T) {
	vs := vmomiTypes.VstorageObjectVCenterQueryChangedDiskAreasRequestType{
		This: vmomiTypes.ManagedObjectReference{
			Type:  "type-1",
			Value: "value-1",
		},
	}

	ns := &v1.VstorageObjectVCenterQueryChangedDiskAreas{}

	if err := v1.FromGovmomi(vs, ns); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, vs.This.Type, ns.This.Type)
}

// TestEnum tests the case where the field is an enum in the govmomi
func TestEnum(t *testing.T) {
	vs := vmomiTypes.PropertyChange{
		Name: "name-1",
		Op:   vmomiTypes.PropertyChangeOpAssign,
	}

	ns := &v1.PropertyChange{}

	if err := v1.FromGovmomi(vs, ns); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, vs.Name, ns.Name)
	assert.Equal(t, vs.Op, *ns.Op.MustToGovmomi())
}

// TestAny tests the case where the field is an any in the govmomi
func TestAny(t *testing.T) {
	vs := vmomiTypes.PropertyChange{
		Name: "name-1",
		Val:  v1.NewPointer("val-1"),
	}

	ns := &v1.PropertyChange{}

	if err := v1.FromGovmomi(vs, ns); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, vs.Name, ns.Name)
}

// TestTime tests the case where the field is a time in the govmomi, and a reminder for the lack of timezone support in protobuf
func TestTime(t *testing.T) {
	vs := vmomiTypes.ScheduledTaskInfo{
		LastModifiedTime: time.Now().UTC(),
	}

	ns := &v1.ScheduledTaskInfo{}

	if err := v1.FromGovmomi(vs, ns); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, vs.LastModifiedTime,
		ns.LastModifiedTime.AsTime())
}

// TestInterface tests the case where the source govmomi type contains an interface
// NOTE: The specific data (like `Limit` value) will lose due to the lack of interface support in protobuf
func TestInterface(t *testing.T) {
	vs := vmomiTypes.LocalizedMethodFault{
		Fault: &vmomiTypes.VmLimitLicense{
			NotEnoughLicenses: vmomiTypes.NotEnoughLicenses{
				RuntimeFault: vmomiTypes.RuntimeFault{
					MethodFault: vmomiTypes.MethodFault{
						FaultCause: nil,
						FaultMessage: []vmomiTypes.LocalizableMessage{
							{
								Key: "key-1",
								Arg: []vmomiTypes.KeyAnyValue{
									{
										Key:   "arg-1",
										Value: 2,
									},
								},
							},
						},
					},
				},
			},
			Limit: 10,
		},
		LocalizedMessage: "?",
	}

	// we will lose Limit value during the conversion since proto has no interface support
	ns := &v1.LocalizedMethodFault{}

	if err := v1.FromGovmomi(vs, ns); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, vs.Fault.GetMethodFault().FaultMessage[0].Key, ns.GetFault().GetFaultMessage()[0].GetKey())
}
