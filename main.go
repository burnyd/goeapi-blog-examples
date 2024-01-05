package main

import (
	"fmt"

	"github.com/aristanetworks/goeapi"
)
// Connection structure this will hold our credentials and other info about the EOS device
type Conn struct {
	Transport string
	Host      string
	Username  string
	Password  string
	Port      int
	Config    string
}

// Json response structure from a show version
type VersionResp struct {
	ModelName        string  `json:"modelName"`
	InternalVersion  string  `json:"internalVersion"`
	SystemMacAddress string  `json:"systemMacAddress"`
	SerialNumber     string  `json:"serialNumber"`
	MemTotal         int     `json:"memTotal"`
	BootupTimestamp  float64 `json:"bootupTimestamp"`
	MemFree          int     `json:"memFree"`
	Version          string  `json:"version"`
	Architecture     string  `json:"architecture"`
	InternalBuildID  string  `json:"internalBuildId"`
	HardwareRevision string  `json:"hardwareRevision,omitempty"`
}

// Simple function that returns the show version string.
func (s *VersionResp) GetCmd() string {
	return "show version"
}

// Method returns a pointer to type of goeapi.Node and a error but connects to the device.
// https://github.com/aristanetworks/goeapi/blob/v1.0.0/client.go#L59
func (c *Conn) Connect() (*goeapi.Node, error) {
	connect, err := goeapi.Connect(c.Transport, c.Host, c.Username, c.Password, c.Port)
	if err != nil {
		fmt.Println(err)
	}
	return connect, nil
}

func main() {
	// Structure the connection data the way we want to structure it.
	d := Conn{
		Transport: "https",
		Host:      "10.255.111.161",
		Username:  "cvpadmin",
		Password:  "cvp123!",
		Port:      443,
	}
	// Use the connection method to connect to the device.
	// This will return the goeapi.Node and a possible error.
	Connect, err := d.Connect()
	if err != nil {
		fmt.Println(err)
	}
	// Print the running-config as a massive string
	// Since type Node has the RunningConfig method if this field is set then 
	// https://github.com/aristanetworks/goeapi/blob/v1.0.0/client.go#L117 that method will be called.
	RunningConfig := Connect.RunningConfig()
	fmt.Println(RunningConfig + "\n")

	// Run some regular commands get the map[string]string output
	fmt.Println("Running a show version \n")
	commands := []string{"show version"}
	// Enable is also a method as well that takes in map[string]string 
	// https://github.com/aristanetworks/goeapi/blob/v1.0.0/client.go#L322
	conf, err := Connect.Enable(commands)
	if err != nil {
		panic(err)
	}
	// Since this returns a []map[string]string we need to range through the first element which is the response.
	for k, v := range conf[0] {
		fmt.Println(k, v)
	}
	// Print out the result.
	fmt.Print(conf[0])

	// Point to the VersionResp struct which has all the json tags
	Showversion := &VersionResp{}
	// Cal the GetHandle method
	handle, err := Connect.GetHandle("json")
	if err != nil {
		fmt.Println(err)
	}
	// This will add to a new slice of AddCommands to send to the switch.
	handle.AddCommand(Showversion)
	// If it exists handle.Call will append all the AddCommands and then connect to the switch
	if err := handle.Call(); err != nil {
		panic(err)
	}
	// Going to print out values for each. 
	fmt.Printf("\n")
	fmt.Printf("Version           : %s\n", Showversion.Version)
	fmt.Printf("System MAC        : %s\n", Showversion.SystemMacAddress)
	fmt.Printf("Serial Number     : %s\n", Showversion.SerialNumber)
}
