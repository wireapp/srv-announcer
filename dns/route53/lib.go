package route53

import (
	"fmt"
	"net"

	"github.com/aws/aws-sdk-go/aws"
	route53Client "github.com/aws/aws-sdk-go/service/route53"
	log "github.com/sirupsen/logrus"
	dns "github.com/zinfra/srv-announcer/dns"
)

// addDotSuffixIfNeeded appends a dot to an existing record, in case it's not there yet.
// TODO: possibly move to parent
func addDotSuffixIfNeeded(dnsName string) string {
	if dnsName[len(dnsName)-1] != byte('.') {
		return dnsName + "."
	}
	return dnsName
}

// editResourceRecords manages a list of rrsets.
// When add is true, it will add the given SRV record to the existing list
// If it already exists (with the exact same weight/priority/port/target values), we keep it where it is
// Not that order matters, but we don't want to shuffle the list before sending it
// When it's false, it will remove the SRV record if it was there
// Removal will only remove entries with the exact same weight/priority/port/target values
func editResourceRecords(add bool, resourceRecords []*route53Client.ResourceRecord, srv *net.SRV) []*route53Client.ResourceRecord {
	newResourceRecords := []*route53Client.ResourceRecord{}

	recordInSet := false

	// copy from the existing records to the new list, ignore the target itself
	for _, rr := range resourceRecords {
		value := aws.StringValue(rr.Value)
		parsed, err := dns.ParseSRV(value)

		if err != nil {
			log.Warnf("Unable to parse SRV record %s, ignoring", value)
			newResourceRecords = append(newResourceRecords, rr)
			continue
		}

		log.Debugf("At %+v", srv)

		// if we see the record we wanted to edit…
		if *parsed == *srv {
			// take note of it, so we don't append on the end
			log.Debugf("Found %v+", srv)
			recordInSet = true
			// if we wanted to remove it, don't copy over
			if !add {
				log.Infof("Found %v+, removing", srv)
				continue
			}
		}

		// just copy over other records
		log.Debugf("Keeping SRV record %s", value)
		newResourceRecords = append(newResourceRecords, rr)
	}

	// if we wanted to add, and the record wasn't in there previously, add to the end
	if add && !recordInSet {
		// synthesize record name
		newValue := fmt.Sprintf("%d %d %d %s", srv.Priority, srv.Weight, srv.Port, srv.Target)
		log.Debugf("Adding SRV record %s", newValue)
		newResourceRecords = append(newResourceRecords, &route53Client.ResourceRecord{
			Value: aws.String(newValue),
		})
	}
	return newResourceRecords
}

// resourceRecordsDiffer compares two records to check if they're different
func resourceRecordsDiffer(a []*route53Client.ResourceRecord, b []*route53Client.ResourceRecord) bool {
	if len(a) != len(b) {
		return true
	}

	for i, e := range a {
		if *e.Value != *b[i].Value {
			return true
		}
	}

	return false
}
