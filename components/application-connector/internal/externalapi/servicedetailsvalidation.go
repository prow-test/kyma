package externalapi

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/kyma-project/kyma/components/application-connector/internal/apperrors"
)

type ServiceDetailsValidator interface {
	Validate(details ServiceDetails) apperrors.AppError
}

type ServiceDetailsValidatorFunc func(details ServiceDetails) apperrors.AppError

func (f ServiceDetailsValidatorFunc) Validate(details ServiceDetails) apperrors.AppError {
	return f(details)
}

func NewServiceDetailsValidator() ServiceDetailsValidator {
	return ServiceDetailsValidatorFunc(func(details ServiceDetails) apperrors.AppError {
		_, err := govalidator.ValidateStruct(details)
		if err != nil {
			return apperrors.WrongInput("incorrect service definition: %s", err.Error())
		}

		if details.Api == nil && details.Events == nil {
			return apperrors.WrongInput(
				"incorrect service definition: at least one of 'api' and 'events' attributes must be provided")
		}

		apperr := validateApiSpec(details.Api)
		if apperr != nil {
			return apperr
		}

		apperr = validateEventsSpec(details.Events)
		if apperr != nil {
			return apperr
		}

		return nil
	})
}

func validateApiSpec(api *API) apperrors.AppError {
	if api != nil && api.Spec != nil {
		err := validateSpec(api.Spec)
		if err != nil {
			return apperrors.WrongInput("api.spec field should be a json object")
		}
	}

	return nil
}

func validateEventsSpec(events *Events) apperrors.AppError {
	if events != nil && events.Spec != nil {
		err := validateSpec(events.Spec)
		if err != nil {
			return apperrors.WrongInput("events.spec field should be a json object")
		}
	}

	return nil
}

func validateSpec(rawMessage json.RawMessage) error {
	var m map[string]*json.RawMessage
	return json.Unmarshal(rawMessage, &m)
}
