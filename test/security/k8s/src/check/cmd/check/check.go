package main

import (
	"flag"
	"log"

	"check/rancher"
	"check/raw"
	"check/validators/master"
)

var (
	ranchercli = flag.Bool("ranchercli", false, "use rancher utility for accessing cluster nodes")
	rke        = flag.Bool("rke", true, "use RKE cluster definition and ssh for accessing cluster nodes (default)")
)

func main() {
	flag.Parse()
	if *ranchercli && *rke {
		log.Fatal("Not supported.")
	}

	var (
		k8sParams []string
		err       error
	)

	switch {
	case *ranchercli:
		k8sParams, err = rancher.GetK8sParams()
	case *rke:
		k8sParams, err = raw.GetK8sParams()
	default:
		log.Fatal("Missing cluster access method.")
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("IsBasicAuthFileAbsent: %t\n", master.IsBasicAuthFileAbsent(k8sParams))
	log.Printf("IsTokenAuthFileAbsent: %t\n", master.IsTokenAuthFileAbsent(k8sParams))
	log.Printf("IsInsecureAllowAnyTokenAbsent: %t\n", master.IsInsecureAllowAnyTokenAbsent(k8sParams))

	log.Printf("IsAnonymousAuthDisabled: %t\n", master.IsAnonymousAuthDisabled(k8sParams))
	log.Printf("IsInsecurePortUnbound: %t\n", master.IsInsecurePortUnbound(k8sParams))
	log.Printf("IsProfilingDisabled: %t\n", master.IsProfilingDisabled(k8sParams))
	log.Printf("IsRepairMalformedUpdatesDisabled: %t\n", master.IsRepairMalformedUpdatesDisabled(k8sParams))
	log.Printf("IsServiceAccountLookupEnabled: %t\n", master.IsServiceAccountLookupEnabled(k8sParams))

	log.Printf("IsKubeletHTTPSAbsentOrEnabled: %t\n", master.IsKubeletHTTPSAbsentOrEnabled(k8sParams))
	log.Printf("IsInsecureBindAddressAbsentOrLoopback: %t\n", master.IsInsecureBindAddressAbsentOrLoopback(k8sParams))
	log.Printf("IsSecurePortAbsentOrValid: %t\n", master.IsSecurePortAbsentOrValid(k8sParams))
}
