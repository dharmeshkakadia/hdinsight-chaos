package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/buger/jsonparser"
)

func GetAmbariComponentURI(cluster string, service string, component string) string {
	return fmt.Sprintf("https://%s.azurehdinsight.net/api/v1/clusters/%s/services/%s/components/%s", cluster, cluster, service, component)
}

func GetAmbariBaseURI(cluster string) string {
	return fmt.Sprintf("https://%s.azurehdinsight.net/api/v1/clusters/%s", cluster, cluster)
}

func GetAmbariReq(uri string, user string, password string, httpverb string) *http.Request {
	req, _ := http.NewRequest(httpverb, uri, nil)
	req.Header.Set("X-Requested-By", "ambari")
	req.SetBasicAuth(user, password)
	return req
}

func GetAmbariNodeList(clustername string, user string, password string) []string {
	var nodes []string
	client := http.DefaultClient
	nodename := regexp.MustCompile("-.*")
	result, _ := client.Do(GetAmbariReq(GetAmbariBaseURI(clustername)+"/hosts", user, password, "get"))
	body, _ := ioutil.ReadAll(result.Body)
	jsonparser.ArrayEach(body, func(value []byte, _ jsonparser.ValueType, o int, e error) {
		h, _ := jsonparser.GetString(value, "Hosts", "host_name")
		nodes = append(nodes, (nodename.ReplaceAllString(h, "")))
	}, "items")
	return nodes
}
