package miner

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type InstanceMetadata struct {
	InstanceType   string `json:"instance_type,omitempty"`
	Hostname       string `json:"hostname,omitempty"`
	PublicHostname string `json:"public_hostname,omitempty"`
}

var AuthStyleEC2Metadata = &Authenticator{
	Name: "X-AWS-EC2-Metadata-Token",
}

func GetInstanceMetadata() *InstanceMetadata {
	result, err := getInstanceMetadata()
	if err != nil {
		fmt.Println("tor-miner:", err)
		return nil
	}
	return result
}

func getInstanceMetadata() (*InstanceMetadata, error) {
	api := &APIEndpoint{
		URL:       "http://169.254.169.254/latest",
		AuthStyle: AuthStyleEC2Metadata,
	}

	req, err := api.NewRequest("PUT", "/api/token", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(api.AuthStyle.Name+"-TTL-Seconds", "21600")

	httpClient := http.Client{Timeout: 1 * time.Second}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("tor-miner: metadata access token:", res.Status)
	} else {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		api.AccessToken = string(body)
	}

	im := &InstanceMetadata{}
	im.Populate(api)
	if im.IsZero() {
		return nil, nil
	}

	return im, nil
}

func (im *InstanceMetadata) Populate(api *APIEndpoint) {
	v := im.reflectedValue()
	for _, tf := range reflect.VisibleFields(v.Type()) {
		key := tf.Tag.Get("json")
		if key == "" {
			continue
		}
		key, _, _ = strings.Cut(key, ",")
		key = strings.ReplaceAll(key, "_", "-")

		res, err := api.Get("/meta-data/" + key)
		if err != nil {
			fmt.Println("tor-miner:", err)
			continue
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			fmt.Printf("tor-miner: metadata %s: %s\n", key, res.Status)
			continue
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("tor-miner:", err)
			continue
		}

		vf := v.FieldByIndex(tf.Index)
		vf.SetString(string(body))
	}
}

func (im *InstanceMetadata) IsZero() bool {
	return im.reflectedValue().IsZero()
}

func (im *InstanceMetadata) reflectedValue() reflect.Value {
	return reflect.Indirect(reflect.ValueOf(im))
}
