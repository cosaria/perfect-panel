package ip

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
	"github.com/pkg/errors"
)

// GetIP parses the input domain or IP address and returns the IP address.
func GetIP(input string) ([]string, error) {
	// Check if the input is already a valid IP address.
	if net.ParseIP(input) != nil {
		return []string{input}, nil
	}

	// Use net.LookupIP to resolve the domain name.
	ips, err := net.LookupIP(input)
	if err != nil {
		return nil, err
	}

	// Convert IP addresses to string format.
	var result []string
	for _, ip := range ips {
		result = append(result, ip.String())
	}
	return result, nil
}

const (
	ipapi   = "ipapi.co"
	ipbase  = "api.ipbase.com"
	ipwhois = "ipwhois.app"
	ipinfo  = "ipinfo.io"
)

var (
	queryUrls = map[string]bool{
		ipbase:  true,
		ipapi:   true,
		ipwhois: true,
		ipinfo:  true,
	}
)

// GetRegionByIp queries the geolocation of an IP address using supported services.
func GetRegionByIp(ip string) (*GeoLocationResponse, error) {
	for service, enabled := range queryUrls {
		if enabled {
			response, err := fetchGeolocation(service, ip)
			if err != nil {
				zap.S().Errorf("Failed to fetch geolocation from %s: %v", service, err)
				continue
			}
			return response, nil
		}
	}
	return nil, errors.New("unable to fetch geolocation for the IP")
}

// fetchGeolocation sends a request to the specified service to retrieve geolocation data.
func fetchGeolocation(service, ip string) (_ *GeoLocationResponse, err error) {
	var apiURL string

	// Construct the API URL based on the service.
	switch service {
	case ipinfo:
		apiURL = fmt.Sprintf("https://ipinfo.io/%s/json", ip)
	case ipapi:
		apiURL = fmt.Sprintf("https://ipapi.co/%s/json", ip)
	case ipbase:
		apiURL = fmt.Sprintf("https://api.ipbase.com/v1/json/%s", ip)
	case ipwhois:
		apiURL = fmt.Sprintf("https://ipwhois.app/json/%s", ip)
	default:
		return nil, fmt.Errorf("unsupported service: %s", service)
	}

	// Create the HTTP request.
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	setHeaders(req, service)

	// Create the HTTP client and send the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close response body: %v", closeErr)
		}
	}()

	// Decompress the response body based on Content-Encoding.
	body, err := decompressResponse(resp)
	if err != nil {
		return nil, err
	}
	// Parse the JSON response into GeoLocationResponse.
	var location GeoLocationResponse
	if err := json.Unmarshal(body, &location); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// Ensure compatibility between country fields.
	if location.Country == "" {
		location.Country = location.CountryName
	}

	if location.Loc != "" && strings.Contains(location.Loc, ",") {
		loc := strings.Split(location.Loc, ",")
		location.Latitude = loc[0]
		location.Longitude = loc[1]
	}

	return &location, nil
}

// setHeaders sets the necessary headers for the HTTP request.
func setHeaders(req *http.Request, host string) {
	req.Header.Set("Host", host)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:134.0) Gecko/20100101 Firefox/134.0")
	req.Header.Set("Accept", "application/json, text/html, application/xhtml+xml, */*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
}

// decompressResponse decompresses the HTTP response body based on its Content-Encoding.
func decompressResponse(resp *http.Response) (_ []byte, err error) {
	var reader io.Reader
	var closeReader func() error

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		gzipReader, gzipErr := gzip.NewReader(resp.Body)
		if gzipErr != nil {
			return nil, fmt.Errorf("failed to create reader: %v", gzipErr)
		}
		reader = gzipReader
		closeReader = gzipReader.Close
	case "br":
		reader = brotli.NewReader(resp.Body)
	case "zstd":
		decoder, zstdErr := zstd.NewReader(resp.Body)
		if zstdErr != nil {
			return nil, fmt.Errorf("failed to create zstd decoder: %v", zstdErr)
		}
		defer decoder.Close()
		return io.ReadAll(decoder)
	default:
		reader = resp.Body
	}

	if closeReader != nil {
		defer func() {
			if closeErr := closeReader(); closeErr != nil && err == nil {
				err = fmt.Errorf("failed to close reader: %v", closeErr)
			}
		}()
	}

	return io.ReadAll(reader)
}

// GeoLocationResponse represents the geolocation data returned by the API.
type GeoLocationResponse struct {
	Country     string `json:"country"`
	CountryName string `json:"country_name"`
	Region      string `json:"region"`
	City        string `json:"city"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	Loc         string `json:"loc"`
}
