package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/cloudflare/cloudflare-go"
	"github.com/rdegges/go-ipify"
)

type Config struct {
	FQDN     string
	APIToken string
}

var config Config

func updateDNSRecord(fqdn string, recordType string, ip string, checkMode bool) error {
	api, err := cloudflare.NewWithAPIToken(config.APIToken)
	if err != nil {
		return err
	}

	domain := domainutil.Domain(fqdn)
	zoneID, err := api.ZoneIDByName(domain)
	if err != nil {
		return err
	}

	fqdnRecord := cloudflare.DNSRecord{Name: fqdn}
	records, err := api.DNSRecords(context.Background(), zoneID, fqdnRecord)
	if err != nil {
		return err
	}

	for _, r := range records {
		if r.Content == ip {
			log.Printf("%s points to current IP address, no change is needed.", r.Name)
			continue
		}
		log.Printf("%s points to %s, the record will be updated.", r.Name, r.Content)

		if checkMode {
			log.Println("Check mode is active, no changes will be made.")
		} else {
			log.Println("Setting", fqdn, "=>", ip, "...")
			r.Content = ip
			err := api.UpdateDNSRecord(context.Background(), zoneID, r.ID, r)
			if err != nil {
				return err
			}
			log.Println("Success!")
		}

	}
	return nil
}

func loadConfig(fn string) error {
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&config)
}

func main() {
	cf := flag.String("c", "config.json", "Config file")
	check_mode := flag.Bool("n", false, "Check mode (dry run)")
	flag.Parse()
	err := loadConfig(*cf)
	if err != nil {
		log.Fatal(err)
	}

	ipv4, err := ipify.GetIp()
	if err != nil {
		log.Fatal("Couldn’t determine the current IP address:", err)
	} else {
		log.Println("Current IP address is:", ipv4)
	}

	err = updateDNSRecord(config.FQDN, "A", ipv4, *check_mode)
	if err != nil {
		log.Println(err)
	}
}
