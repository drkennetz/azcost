package azure

import (
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/costmanagement/armcostmanagement"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// NewResourceIdTypeGroupGrouping returns a new QueryGrouping querying ResourceId, ResourceType, and ResourceGroup
func NewResourceIdTypeGroupGrouping(queryColumnType string) []*armcostmanagement.QueryGrouping {
	var grouping []*armcostmanagement.QueryGrouping
	newGroupingId := NewQueryGrouping(QueryGroupingResourceId, queryColumnType)
	newGroupingType := NewQueryGrouping(QueryGroupingResourceType, queryColumnType)
	newGroupingGroup := NewQueryGrouping(QueryGroupingResourceGroup, queryColumnType)
	grouping = append(grouping, &newGroupingId, &newGroupingType, &newGroupingGroup)
	return grouping
}

// NewResourceTypeGroupGrouping returns a new QueryGrouping querying ResourceType and ResourceGroup
func NewResourceTypeGroupGrouping(queryColumnType string) []*armcostmanagement.QueryGrouping {
	var grouping []*armcostmanagement.QueryGrouping
	newGroupingType := NewQueryGrouping(QueryGroupingResourceType, queryColumnType)
	newGroupingGroup := NewQueryGrouping(QueryGroupingResourceGroup, queryColumnType)
	grouping = append(grouping, &newGroupingType, &newGroupingGroup)
	return grouping
}

func GetRequest() {
	costUrl := "https://management.azure.com/subscriptions/475b3c71-8110-4997-9e77-2f86d6c6cd42/providers/Microsoft.CostManagement/dimensions?api-version=2021-10-01&$expand=properties/data&$skiptoken=AQAAAA%3D%3D"
	authendpoint := "https://login.microsoftonline.com/22340fa8-9226-4871-b677-d3b3e377af72/oauth2/token"
	body := url.Values(map[string][]string{
		"resource":      {"https://management.azure.com"},
		"client_id":     {"7064d835-fe3c-435f-aa50-837db5a3c60b"},
		"client_secret": {"Ixf8Q~Q-J8y8idaUjX-CWzTy~VXt8kf3Kip0KaEW"},
		"grant_type":    {"client_credentials"}})

	request, err := http.NewRequest(
		http.MethodPost,
		authendpoint,
		strings.NewReader(body.Encode()))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	var loginResp loginResponse
	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		if err = getJson(resp, &loginResp); err != nil {
			log.Fatalln(err)
		}
	}

	var tok = "Bearer " + loginResp.Token
	costReq, err := http.NewRequest("POST", costUrl, nil)
	costReq.Header.Add("Authorization", tok)
	costClient := &http.Client{}
	costResp, err := costClient.Do(costReq)
	if err != nil {
		log.Fatalln(err)
	}
	defer costResp.Body.Close()
	costBody, err := io.ReadAll(costResp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(costBody))
}

func getJson(resp *http.Response, target *loginResponse) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

type loginResponse struct {
	Token        string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	ExpiresOn    string `json:"expires_on"`
	ExtExpiresIn string `json:"ext_expires_in"`
	NotBefore    string `json:"not_before"`
	Resource     string `json:"resource"`
	TokenType    string `json:"token_type"`
}
