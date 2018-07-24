package v2

import (
	istioAuthApi "github.com/kyma-project/kyma/components/api-controller/pkg/apis/authentication.istio.io/v1alpha1"
	kymaMeta "github.com/kyma-project/kyma/components/api-controller/pkg/apis/gateway.kyma.cx/meta/v1"
	istioAuth "github.com/kyma-project/kyma/components/api-controller/pkg/clients/authentication.istio.io/clientset/versioned"
	istioAuthTyped "github.com/kyma-project/kyma/components/api-controller/pkg/clients/authentication.istio.io/clientset/versioned/typed/authentication.istio.io/v1alpha1"
	"github.com/kyma-project/kyma/components/api-controller/pkg/controller/commons"
	"github.com/kyma-project/kyma/components/api-controller/pkg/controller/meta"
	log "github.com/sirupsen/logrus"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	k8sMeta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
)

type istioImpl struct {
	istioAuthInterface istioAuth.Interface
}

func New(a istioAuth.Interface) Interface {
	return &istioImpl{
		istioAuthInterface: a,
	}
}

func (a *istioImpl) Create(dto *Dto) (*kymaMeta.GatewayResource, error) {

	if isAuthenticationDisabled(dto) {
		return nil, nil
	}

	istioAuthPolicy := toIstioAuthPolicy(dto)

	log.Debugf("Creating authentication policy: %v", istioAuthPolicy)

	created, err := a.istioAuthPolicyInterface(dto.MetaDto).Create(istioAuthPolicy)
	if err != nil {
		return nil, commons.HandleError(err, "error while creating authentication policy")
	}

	log.Debugf("Authentication policy created: %v", istioAuthPolicy)
	return gatewayResourceFrom(created), nil
}

func (a *istioImpl) Update(oldDto, newDto *Dto) (*kymaMeta.GatewayResource, error) {

	if isAuthenticationDisabled(newDto) {

		log.Debugf("Authentication disabled. Trying to delete the old authentication policy...")
		// no new newRule; we only have to delete the old one, if it exists
		if err := a.Delete(oldDto); err != nil {
			return nil, err
		}
		return nil, nil
	}

	// there is a authentication policy to update / create
	newIstioAuthPolicy := toIstioAuthPolicy(newDto)

	log.Debugf("Authentication enabled. Trying to create or update authentication policy: %v", newIstioAuthPolicy)

	// checking if authentication policy has to be created
	if oldDto == nil || oldDto.Status.Resource.Name == "" {

		log.Debug("Authentication policy does not exist. Creating...")

		// create newRule
		createdResource, err := a.Create(newDto)

		if err != nil {
			return nil, commons.HandleError(err, "error while recreating authentication policy (can not create a new one)")
		}

		log.Debugf("Authentication policy created: %v", createdResource)
		return createdResource, nil
	}

	oldIstioAuthPolicy := toIstioAuthPolicy(oldDto)

	if a.isEqual(oldIstioAuthPolicy, newIstioAuthPolicy) {

		log.Debugf("Update skipped: authentication policy has not changed.")
		return nil, nil
	}

	newIstioAuthPolicy.ObjectMeta.ResourceVersion = oldDto.Status.Resource.Version

	// newRule should be updated (i was not recreated and it was differs from the old one)
	log.Debugf("Updating authentication policy: %v", newIstioAuthPolicy)

	updated, err := a.istioAuthPolicyInterface(newDto.MetaDto).Update(newIstioAuthPolicy)
	if err != nil {
		return nil, commons.HandleError(err, "error while updating authentication policy")
	}

	log.Debugf("Authentication policy updated: %v", updated)
	return gatewayResourceFrom(updated), nil
}

func (a *istioImpl) Delete(dto *Dto) error {

	if dto == nil {
		log.Debug("Delete skipped: no authentication policy to delete.")
		return nil
	}
	return a.deleteByName(dto.MetaDto)
}

func (a *istioImpl) deleteByName(meta meta.Dto) error {

	// if there is no rule to delete, just skip it
	if meta.Name == "" {
		log.Debug("Delete skipped: no authentication policy to delete.")
		return nil
	}
	log.Debugf("Deleting authentication policy: %s", meta.Name)

	err := a.istioAuthPolicyInterface(meta).Delete(meta.Name, &k8sMeta.DeleteOptions{})
	if err != nil && !apiErrors.IsNotFound(err) {
		return commons.HandleError(err, "error while deleting authentication policy")
	}

	log.Debugf("Authentication policy deleted: %+v", meta.Name)
	return nil
}

func (a *istioImpl) istioAuthPolicyInterface(metaDto meta.Dto) istioAuthTyped.PolicyInterface {
	return a.istioAuthInterface.AuthenticationV1alpha1().Policies(metaDto.Namespace)
}

func (a *istioImpl) isEqual(oldRule *istioAuthApi.Policy, newRule *istioAuthApi.Policy) bool {
	return reflect.DeepEqual(oldRule.Spec, newRule.Spec)
}

func toIstioAuthPolicy(dto *Dto) *istioAuthApi.Policy {

	objectMetadata := k8sMeta.ObjectMeta{
		Name:      dto.MetaDto.Name,
		Namespace: dto.MetaDto.Namespace,
		Labels:    dto.MetaDto.Labels,
	}

	spec := &istioAuthApi.PolicySpec{
		Targets: []*istioAuthApi.Target{
			{Name: dto.ServiceName},
		},
	}

	origins := make(istioAuthApi.Origins, 0, 1)
	for _, rule := range dto.Rules {

		if rule.Type == JwtType {
			origins = append(origins, &istioAuthApi.Origin{
				Jwt: &istioAuthApi.Jwt{
					Issuer:  rule.Jwt.Issuer,
					JwksUri: rule.Jwt.JwksUri,
				},
			})
		}
	}
	spec.Origins = origins

	spec.PrincipalBinding = istioAuthApi.UseOrigin

	return &istioAuthApi.Policy{
		ObjectMeta: objectMetadata,
		Spec:       spec,
	}
}

func isAuthenticationDisabled(dto *Dto) bool {
	return dto == nil || len(dto.Rules) == 0
}

func gatewayResourceFrom(policy *istioAuthApi.Policy) *kymaMeta.GatewayResource {
	return &kymaMeta.GatewayResource{
		Name:    policy.Name,
		Uid:     policy.UID,
		Version: policy.ResourceVersion,
	}
}
