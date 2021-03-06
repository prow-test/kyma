package externalapi

import (
	"encoding/json"
	"github.com/kyma-project/kyma/components/application-connector/internal/apperrors"
	"github.com/kyma-project/kyma/components/application-connector/internal/metadata"
	"github.com/kyma-project/kyma/components/application-connector/internal/metadata/serviceapi"
)

type Service struct {
	ID          string `json:"id"`
	Provider    string `json:"provider"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ServiceDetails struct {
	Provider      string         `json:"provider" valid:"required~Provider field cannot be empty."`
	Name          string         `json:"name" valid:"required~Name field cannot be empty."`
	Description   string         `json:"description" valid:"required~Description field cannot be empty."`
	Api           *API           `json:"api,omitempty"`
	Events        *Events        `json:"events,omitempty"`
	Documentation *Documentation `json:"documentation,omitempty"`
}

type CreateServiceResponse struct {
	ID string `json:"id"`
}

type API struct {
	TargetUrl   string          `json:"targetUrl" valid:"url,required~targetUrl field cannot be empty."`
	Credentials *Credentials    `json:"credentials,omitempty"`
	Spec        json.RawMessage `json:"spec,omitempty"`
}

type Credentials struct {
	Oauth Oauth `json:"oauth" valid:"required~oauth field cannot be empty"`
}

type Oauth struct {
	URL          string `json:"url" valid:"url,required~oauth url field cannot be empty"`
	ClientID     string `json:"clientId" valid:"required~oauth clientId field cannot be empty"`
	ClientSecret string `json:"clientSecret" valid:"required~oauth clientSecret cannot be empty"`
}

type Events struct {
	Spec json.RawMessage `json:"spec" valid:"required~spec cannot be empty"`
}

type Documentation struct {
	DisplayName string       `json:"displayName" valid:"required~displayName field cannot be empty in documentation"`
	Description string       `json:"description" valid:"required~description field cannot be empty in documentation"`
	Type        string       `json:"type" valid:"required~type field cannot be empty in documentation"`
	Tags        []string     `json:"tags,omitempty"`
	Docs        []DocsObject `json:"docs,omitempty"`
}

type DocsObject struct {
	Title  string `json:"title"`
	Type   string `json:"type"`
	Source string `json:"source"`
}

const stars = "********"

func serviceDefinitionToService(serviceDefinition metadata.ServiceDefinition) Service {
	return Service{
		ID:          serviceDefinition.ID,
		Name:        serviceDefinition.Name,
		Provider:    serviceDefinition.Provider,
		Description: serviceDefinition.Description,
	}
}

func serviceDefinitionToServiceDetails(serviceDefinition metadata.ServiceDefinition) (ServiceDetails, apperrors.AppError) {
	serviceDetails := ServiceDetails{
		Provider:    serviceDefinition.Provider,
		Name:        serviceDefinition.Name,
		Description: serviceDefinition.Description,
	}

	if serviceDefinition.Api != nil {
		serviceDetails.Api = &API{
			TargetUrl: serviceDefinition.Api.TargetUrl,
			Spec:      serviceDefinition.Api.Spec,
		}

		if serviceDefinition.Api.Credentials != nil {
			serviceDetails.Api.Credentials = &Credentials{
				Oauth: Oauth{
					ClientID:     stars,
					ClientSecret: stars,
					URL:          serviceDefinition.Api.Credentials.Oauth.URL,
				},
			}
		}
	}

	if serviceDefinition.Events != nil {
		serviceDetails.Events = &Events{
			Spec: serviceDefinition.Events.Spec,
		}
	}

	if serviceDefinition.Documentation != nil {
		err := json.Unmarshal(serviceDefinition.Documentation, &serviceDetails.Documentation)
		if err != nil {
			return serviceDetails, apperrors.Internal("failed to unmarshal documentation, '%s'", err)
		}

	}

	return serviceDetails, nil
}

func serviceDetailsToServiceDefinition(serviceDetails ServiceDetails) (metadata.ServiceDefinition, apperrors.AppError) {
	serviceDefinition := metadata.ServiceDefinition{
		Provider:    serviceDetails.Provider,
		Name:        serviceDetails.Name,
		Description: serviceDetails.Description,
	}

	if serviceDetails.Api != nil {
		serviceDefinition.Api = &serviceapi.API{
			TargetUrl: serviceDetails.Api.TargetUrl,
		}
		if serviceDetails.Api.Credentials != nil {
			serviceDefinition.Api.Credentials = &serviceapi.Credentials{
				Oauth: serviceapi.Oauth{
					ClientID:     serviceDetails.Api.Credentials.Oauth.ClientID,
					ClientSecret: serviceDetails.Api.Credentials.Oauth.ClientSecret,
					URL:          serviceDetails.Api.Credentials.Oauth.URL,
				},
			}
		}
		if serviceDetails.Api.Spec != nil {
			serviceDefinition.Api.Spec = compact(serviceDetails.Api.Spec)
		}
	}

	if serviceDetails.Events != nil && serviceDetails.Events.Spec != nil {
		serviceDefinition.Events = &metadata.Events{
			Spec: compact(serviceDetails.Events.Spec),
		}
	}

	if serviceDetails.Documentation != nil {
		marshalled, err := json.Marshal(&serviceDetails.Documentation)
		if err != nil {
			return serviceDefinition, apperrors.WrongInput("failed to marshal documentation, '%s'", err)
		}
		serviceDefinition.Documentation = marshalled
	}

	return serviceDefinition, nil
}
