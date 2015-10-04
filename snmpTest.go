package main

import (
	"flag"
	"fmt"
	"github.com/alouca/gosnmp"
)

var (
	cmdCommunity string
	cmdTarget    string
	cmdOid       string
	cmdDebug     string
	cmdTimeout   int64
)

func init() {
  fmt.Printf("Simple SNMP Query ( just a simple GO exercise )\n")
	flag.StringVar(&cmdDebug, "debug", "", "Debug flag expects byte array of raw packet to test decoding")

	flag.StringVar(&cmdTarget, "target", "", "Target SNMP Agent")
	flag.StringVar(&cmdCommunity, "community", "", "SNNP Community")
	flag.StringVar(&cmdOid, "oid", ".1.3.6.1.2.1.1.1.0", "OID")
	flag.Int64Var(&cmdTimeout, "timeout", 5, "Set the timeout in seconds")
	flag.Parse()
}

func debugOut(msg string) int {
    if ( cmdDebug != "" ) {
        fmt.Printf("DEBUG --> %s\n", msg)
    }
    return 1;
}

func main() {
    debugOut("Running in DEBUG mode")
    s, err := gosnmp.NewGoSNMP(cmdTarget, cmdCommunity, gosnmp.Version2c, 5)
	if cmdDebug != "" {
        debugOut(fmt.Sprintf("Sending SNMP Bulk Query against %s on OID %s via community %s.", cmdTarget, cmdOid, cmdCommunity))
		s.SetDebug(false)
		s.SetVerbose(true)
	} else {
        s.SetDebug(false)
        s.SetVerbose(false)
	}

	if cmdTarget == "" {
		flag.PrintDefaults()
		return
	}

	if err != nil {
		fmt.Printf("Error creating SNMP instance: %s\n", err.Error())
		return
	}

	s.SetTimeout(cmdTimeout)
	fmt.Printf("Getting %s\n", cmdOid)
	resp, err := s.Get(cmdOid)
	if err != nil {
		fmt.Printf("Error getting response: %s\n", err.Error())
	} else {
		for _, v := range resp.Variables {
			fmt.Printf("%s -> ", v.Name)
			switch v.Type {
			case gosnmp.OctetString:
				if s, ok := v.Value.(string); ok {
					fmt.Printf("%s\n", s)
				} else {
					fmt.Printf("Response is not a string\n")
				}
			default:
				fmt.Printf("Type: %d - Value: %v\n", v.Type, v.Value)
			}
		}

	}
}
