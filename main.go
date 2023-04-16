package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const urlPrefix string = "https://api.netlify.com/api/v1/"

type DnsZone struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type DnsRecord struct {
	Id        string  `json:"id"`
	DnsZoneId string  `json:"dns_zone_id"`
	Hostname  string  `json:"hostname"`
	Type      string  `json:"type"`
	Ttl       int     `json:"ttl"`
	Priority  int     `json:"priority"`
	Weight    *int    `json:"weight,omitempty"`
	Port      *int    `json:"port,omitempty"`
	Flag      *string `json:"flag,omitempty"`
	Tag       *string `json:"tag,omitempty"`
	Managed   bool    `json:"managed"`
	Value     string  `json:"value"`
}

type NetlifyDnsClient struct {
	client *http.Client
	token  string
}

func NewNetlifyDnsClient(token string) NetlifyDnsClient {
	client := &http.Client{}

	return NetlifyDnsClient{client: client, token: token}
}

func (n *NetlifyDnsClient) addAuthHeader(req *http.Request) {
	req.Header.Add("Authorization", "Bearer "+n.token)
}

func (n *NetlifyDnsClient) getReqByteSlice(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", urlPrefix+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating get request: %w", err)
	}
	n.addAuthHeader(req)

	resp, err := n.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error doing get request: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading get request body: %w", err)
	}

	return body, nil
}

func (n *NetlifyDnsClient) GetAllDnsZones() ([]DnsZone, error) {
	body, err := n.getReqByteSlice("dns_zones")
	if err != nil {
		return nil, err
	}

	var dnsZones []DnsZone
	err = json.Unmarshal(body, &dnsZones)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling get request body: %w", err)
	}

	return dnsZones, nil
}

func (n *NetlifyDnsClient) GetAllDnsRecords(zoneId string) ([]DnsRecord, error) {
	body, err := n.getReqByteSlice("dns_zones/" + zoneId + "/dns_records")
	if err != nil {
		return nil, err
	}

	var records []DnsRecord
	err = json.Unmarshal(body, &records)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling get request body: %w", err)
	}

	return records, nil
}

func GenerateZoneFile(zone DnsZone, records []DnsRecord) (string, error) {
	var zoneFile strings.Builder

	zoneFile.WriteString(fmt.Sprintf("$ORIGIN %s\n", zone.Name+"."))

	for _, record := range records {
		if record.Type == "NETLIFY" {
			fmt.Printf("ingoring NETLIFY record: %s\n", record.Hostname)
			continue
		}

		name := record.Hostname + "."

		var value string
		if record.Type == "CNAME" {
			value = record.Value + "."
		} else {
			value = record.Value
		}

		zoneFile.WriteString(
			fmt.Sprintf(
				"%s\tIN\t%d\t%s\t%s\n",
				name,
				record.Ttl,
				record.Type,
				value,
			),
		)
	}

	return zoneFile.String(), nil
}

func main() {
	token := os.Getenv("NETLIFY_TOKEN")

	client := NewNetlifyDnsClient(token)

	zones, err := client.GetAllDnsZones()
	if err != nil {
		log.Fatalln(err)
	}
	for _, zone := range zones {
		records, err := client.GetAllDnsRecords(zone.Id)
		if err != nil {
			log.Fatalln(err)
		}

		zoneContents, err := GenerateZoneFile(zone, records)
		if err != nil {
			log.Fatalln(err)
		}

		fileName := zone.Id+".zone"

		err = os.WriteFile(fileName, []byte(zoneContents), 0644)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(fileName)
	}
}
