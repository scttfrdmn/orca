package aws

import "time"

// Instance represents an EC2 instance.
type Instance struct {
	ID           string
	Type         string
	State        string
	PublicIP     string
	PrivateIP    string
	LaunchTime   time.Time
	InstanceType string
}
