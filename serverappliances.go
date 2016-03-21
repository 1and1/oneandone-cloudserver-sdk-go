package oneandone

import "net/http"

type ServerAppliance struct {
	Identity
	typeField
	OsImageType        string      `json:"os_image_type,omitempty"`
	OsFamily           string      `json:"os_family,omitempty"`
	Os                 string      `json:"os,omitempty"`
	OsVersion          string      `json:"os_version,omitempty"`
	MinHddSize         int         `json:"min_hdd_size"`
	Architecture       interface{} `json:"architecture"`
	Licenses           []License   `json:"licenses,omitempty"`
	IsAutomaticInstall bool        `json:"automatic_installation"`
	AvailableSites     []string    `json:"available_sites,omitempty"`
	ApiPtr
}

// GET /server_appliances
func (api *API) ListServerAppliances(args ...interface{}) ([]ServerAppliance, error) {
	url, err := processQueryParams(createUrl(api, serverAppliancePathSegment), args...)
	if err != nil {
		return nil, err
	}
	res := []ServerAppliance{}
	err = api.Client.Get(url, &res, http.StatusOK)
	if err != nil {
		return nil, err
	}
	for index, _ := range res {
		res[index].api = api
	}
	return res, nil
}

// GET /server_appliances/{id}
func (api *API) GetServerAppliance(sa_id string) (*ServerAppliance, error) {
	res := new(ServerAppliance)
	url := createUrl(api, serverAppliancePathSegment, sa_id)
	err := api.Client.Get(url, &res, http.StatusOK)
	if err != nil {
		return nil, err
	}
	//	res.api = api
	return res, nil
}
