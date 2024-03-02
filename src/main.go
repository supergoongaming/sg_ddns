package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

const (
	zonesEnvName  = "DNS_ZONES"
	zoneIdEnvName = "ZONE_ID"
)

func getPublicIP() string {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		log.Fatalf("Could not get public IP address, will not update.  Error: %s", err)
	}
	defer resp.Body.Close()
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Could not get public IP address, will not update.  Error: %s", err)
	}
	return string(ip)
}

func getAllChanges(zones string, ip string) ([]*route53.Change, error) {
	var changes []*route53.Change
	for _, zone := range strings.Split(zones, " ") {
		recordName := fmt.Sprintf("%s.supergoon.com", zone)
		fmt.Printf("Going to update record %s", recordName)
		newARecord := &route53.ResourceRecordSet{
			Name:            &recordName,
			Type:            aws.String("A"),
			TTL:             aws.Int64(300),
			ResourceRecords: []*route53.ResourceRecord{{Value: aws.String(ip)}},
		}
		change := &route53.Change{
			Action:            aws.String("UPSERT"),
			ResourceRecordSet: newARecord,
		}
		changes = append(changes, change)
	}
	if len(changes) == 0 {
		return nil, fmt.Errorf("couldn't create any changes for zones %s for ip %s", zones, ip)
	}
	return changes, nil

}

func main() {
	newIp := getPublicIP()
	zones := os.Getenv(zonesEnvName)
	hostedZoneId := os.Getenv(zoneIdEnvName)
	if len(zones) == 0 || len(hostedZoneId) == 0 {
		log.Fatalf("Error getting environment variables! These should not be empty. Key:Value\n%s is %s\n%s is %s", zonesEnvName, zones, zoneIdEnvName, hostedZoneId)
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	r53Client := route53.New(sess)
	changes, err := getAllChanges(zones, newIp)
	if err != nil {
		log.Fatalf("Error getting changes\nError: %s", err)
	}
	changeBatch := &route53.ChangeBatch{
		Changes: changes,
	}
	changeParams := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(hostedZoneId),
		ChangeBatch:  changeBatch,
	}
	_, err = r53Client.ChangeResourceRecordSets(changeParams)
	if err != nil {
		log.Fatalf("Failed to update A record: %v", err)
	}
	fmt.Println("A record updated successfully.")
}
