package ccv2

import (
	"encoding/json"

	"code.cloudfoundry.org/cli/api/cloudcontroller"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv2/internal"
)

type SecurityGroupRule struct {
	Destination string
	Ports       string
	Protocol    string
}

type SecurityGroup struct {
	Description string
	GUID        string
	Name        string
	Rules       []SecurityGroupRule
}

// UnmarshalJSON helps unmarshal a Cloud Controller Security Group response
func (securityGroup *SecurityGroup) UnmarshalJSON(data []byte) error {
	var ccSecurityGroup struct {
		Metadata internal.Metadata `json:"metadata"`
		Entity   struct {
			GUID string `json:"guid"`
			Name string `json:"name"`
			// Rules []struct {
			// 	Destination string `json:"destination"`
			// 	Ports       string `json:"ports"`
			// 	Protocol    string `json:"protocol"`
			// } `json:"rules"`
		} `json:"entity"`
	}

	if err := json.Unmarshal(data, &ccSecurityGroup); err != nil {
		return err
	}

	securityGroup.GUID = ccSecurityGroup.Metadata.GUID
	securityGroup.Name = ccSecurityGroup.Entity.Name
	// securityGroup.Rules = make([]SecurityGroupRule, len(ccSecurityGroup.Entity.Rules))
	// for i, ccRule := range ccSecurityGroup.Entity.Rules {
	// 	securityGroup.Rules[i].Destination = ccRule.Destination
	// 	securityGroup.Rules[i].Ports = ccRule.Ports
	// 	securityGroup.Rules[i].Protocol = ccRule.Protocol
	// }
	return nil
}

func (client *Client) AssociateSpaceWithSecurityGroup(securityGroupGUID string, spaceGUID string) (Warnings, error) {
	request, err := client.newHTTPRequest(requestOptions{
		RequestName: internal.AssociateSpaceWithSecurityGroupRequest,
		URIParams: Params{
			"security_group_guid": securityGroupGUID,
			"space_guid":          spaceGUID,
		},
	})

	if err != nil {
		return nil, err
	}

	response := cloudcontroller.Response{}

	err = client.connection.Make(request, &response)
	return response.Warnings, err
}

func (client *Client) GetSecurityGroups(queries []Query) ([]SecurityGroup, Warnings, error) {
	request, err := client.newHTTPRequest(requestOptions{
		RequestName: internal.SecurityGroupsRequest,
		Query:       FormatQueryParameters(queries),
	})

	if err != nil {
		return nil, nil, err
	}

	var securityGroupsList []SecurityGroup
	warnings, err := client.paginate(request, SecurityGroup{}, func(item interface{}) error {
		if securityGroup, ok := item.(SecurityGroup); ok {
			securityGroupsList = append(securityGroupsList, securityGroup)
		} else {
			return cloudcontroller.UnknownObjectInListError{
				Expected:   SecurityGroup{},
				Unexpected: item,
			}
		}
		return nil
	})

	return securityGroupsList, warnings, err
}

// GetSharedSecurityGroup returns the Shared Domain associated with the provided
// Space GUID.
func (client *Client) GetSharedSecurityGroup(spaceGUID string) (SecurityGroup, Warnings, error) {
	// request, err := client.newHTTPRequest(requestOptions{
	// 	RequestName: internal.SharedSecurityGroup,
	// 	URIParams:   map[string]string{"shared_domain_guid": domainGUID},
	// })
	// if err != nil {
	// 	return SecurityGroup{}, nil, err
	// }

	// var domain SecurityGroup
	// response := cloudcontroller.Response{
	// 	Result: &domain,
	// }

	// err = client.connection.Make(request, &response)
	// if err != nil {
	// 	return SecurityGroup{}, response.Warnings, err
	// }

	// return domain, response.Warnings, nil
	return SecurityGroup{}, nil, nil
}

func (client *Client) GetSpaceRunningSecurityGroupsBySpace(spaceGUID string) ([]SecurityGroup, Warnings, error) {
	return nil, nil, nil
}

func (client *Client) GetSpaceStagingSecurityGroupsBySpace(spaceGUID string) ([]SecurityGroup, Warnings, error) {
	return nil, nil, nil
}
