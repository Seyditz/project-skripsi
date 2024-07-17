package utils

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Seyditz/project-skripsi/database"
	"github.com/Seyditz/project-skripsi/models"
)

type OAIResponse struct {
	Records         []Record        `xml:"ListRecords>record"`
	ResumptionToken ResumptionToken `xml:"ListRecords>resumptionToken"`
}
type Record struct {
	Metadata Metadata `xml:"metadata>dc"`
}

type Metadata struct {
	Title     string   `xml:"title"`
	Publisher string   `xml:"publisher"`
	Subject   []string `xml:"subject"`
}

type ResumptionToken struct {
	Value string `xml:",chardata"`
}

func FetchTitles() ([]string, error) {
	baseURL := "http://repository.upnvj.ac.id/cgi/oai2?verb=ListRecords&metadataPrefix=oai_dc"
	var titles []string

	for {
		resp, err := http.Get(baseURL)
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil data: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("gagal membaca respons: %v", err)
		}

		var oaiResponse OAIResponse
		err = xml.Unmarshal(body, &oaiResponse)
		if err != nil {
			return nil, fmt.Errorf("tidak dapat parsing XML: %v", err)
		}

		for _, record := range oaiResponse.Records {
			for _, subject := range record.Metadata.Subject {
				if strings.Contains(strings.ToLower(subject), "computer science") {
					titles = append(titles, record.Metadata.Title)
					break
				}
			}
		}

		// break

		if oaiResponse.ResumptionToken.Value == "" {
			break
		}

		baseURL = "http://repository.upnvj.ac.id/cgi/oai2?verb=ListRecords&resumptionToken=" + oaiResponse.ResumptionToken.Value
	}

	return titles, nil
}

func GetRepoTitlesFromDatabase() ([]string, error) {
	var repoJudul []models.Title
	database.DB.Find(&repoJudul)

	if result := database.DB.Find(&repoJudul); result.RowsAffected == 0 {
		return nil, fmt.Errorf("tidak bisa fetching judul")
	}

	juduls := repoJudul[1].Titles

	return juduls, nil
}

var softwareEngineerKeywords = []string{
	"software", "aplikasi", "program", "coding", "pengembangan", "algoritma",
	"pengujian", "sistem", "komputasi", "rekayasa", "perangkat lunak",
	"pemrograman", "aplikasi mobile", "aplikasi web", "framework", "database",
	"react", "angular", "vue", "django", "flask", "spring", "java", "python",
	"javascript", "typescript", "golang", "node.js", "rails", "laravel",
	"android", "ios", "kotlin", "swift", "ci/cd", "docker", "kubernetes",
	"cloud", "aws", "azure", "gcp", "search engine", "ocr", "ansys", "kansei engineering",
	"genetika", "kotlin", "unity", "machine learning", "data mining", "solidworks",
	"cf", "vissim", "hecras", "gis", "biogeme", "augmented reality", "virtual reality",
	"flap", "simulation", "finite element method", "vibration", "modelling",
	"paving block", "repo", "dspace", "winqsb", "robotics", "3d modeling", "spring boot",
}

var itSecuritySpecialistKeywords = []string{
	"keamanan", "security", "cryptography", "enkripsi", "penetration testing",
	"firewall", "vulnerability", "malware", "cyber", "phishing",
	"autentikasi", "otorisasi", "keamanan jaringan", "serangan siber",
	"kaspesky", "norton", "mcafee", "wireshark", "nmap", "metasploit",
	"burpsuite", "nessus", "owasp", "snort", "suricata", "ids", "ips",
	"blockchain", "forensik", "siem", "splunk", "logrhythm", "graylog",
	"ssl", "tls", "vpn", "two-factor authentication", "hacking", "cyber crime",
	"phishing", "password", "brute force", "intrusion detection", "intrusion prevention",
	"biometric", "cloud security", "iot security", "gdpr", "data protection",
	"network security", "firewall", "cryptanalysis", "rsa", "aes", "des",
	"public key infrastructure", "digital forensics", "identity management",
	"encryption", "threat intelligence", "siem", "compliance", "risk management",
	"security policies", "incident response", "zero trust", "honeypot",
	"key management", "data breach", "secure coding", "sql injection",
	"xss", "cross-site scripting", "dos", "denial of service", "siber",
}

var networkEngineerKeywords = []string{
	"sdn", "software-defined networking", "load balancer", "ids", "intrusion detection system",
	"snmp", "vpn", "virtual private network", "bgp", "ospf", "wan", "mesh", "wi-fi 6",
	"kriptografi", "iot", "internet of things", "vlan", "hsrp", "vrrp", "wireshark",
	"qos", "voip", "firewall", "mpls", "fiber optik", "honeypot", "5g", "data center",
	"ipv6", "nac", "eigrp", "nfv", "network function virtualization", "wpa3", "mikrotik",
	"openstack", "site-to-site vpn", "sd-wan", "ftth", "nagios", "palo alto", "netflow",
	"ssl/tls", "ansible", "lte", "wi-fi", "li-fi", "openflow", "vxlan", "mesh network",
	"python", "isp", "openvpn", "gre", "mimo", "suricata", "metro ethernet", "cloud",
	"siem", "splunk", "zabbix", "802.11ac", "backbone", "ipsec", "pfsense", "lisp",
	"snort", "ssl vpn", "chef", "puppet", "10 gigabit ethernet", "ssh", "blockchain",
	"sflow", "40 gigabit ethernet", "tls 1.3", "802.11ad", "is-is", "artificial intelligence",
}

func ClassifyTitle(title string) string {
	titleLower := strings.ToLower(title)

	for _, keyword := range softwareEngineerKeywords {
		if strings.Contains(titleLower, keyword) {
			return "Software Engineer"
		}
	}

	for _, keyword := range itSecuritySpecialistKeywords {
		if strings.Contains(titleLower, keyword) {
			return "IT Security Specialist"
		}
	}

	for _, keyword := range networkEngineerKeywords {
		if strings.Contains(titleLower, keyword) {
			return "Network Engineer"
		}
	}

	return "Unclassified"
}
