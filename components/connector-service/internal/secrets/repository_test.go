package secrets

import (
	"github.com/kyma-project/kyma/components/connector-service/internal/apperrors"
	"github.com/kyma-project/kyma/components/connector-service/internal/secrets/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

const (
	reName = "reName"
)

var (
	expectedCaCrt = []byte("caCrtEncoded")
	expectedCaKey = []byte("caKeyEncoded")
)

func TestRepository_Get(t *testing.T) {

	t.Run("should get secret", func(t *testing.T) {
		// given
		secretMap := make(map[string][]byte)
		secretMap["ca.crt"] = expectedCaCrt
		secretMap["ca.key"] = expectedCaKey

		secretsManager := &mocks.Manager{}
		secretsManager.On("Get", reName, metav1.GetOptions{}).Return(&v1.Secret{Data: secretMap}, nil)

		repository := NewRepository(secretsManager)

		// when
		encodedCrt, encodedKey, err := repository.Get(reName)

		// then
		require.NoError(t, err)

		assert.Equal(t, expectedCaCrt, encodedCrt)
		assert.Equal(t, expectedCaKey, encodedKey)
	})

	t.Run("should fail in case secret not found", func(t *testing.T) {
		// given
		k8sNotFoundError := &k8serrors.StatusError{
			ErrStatus: metav1.Status{Reason: metav1.StatusReasonNotFound},
		}
		secretsManager := &mocks.Manager{}
		secretsManager.On("Get", reName, metav1.GetOptions{}).Return(nil, k8sNotFoundError)

		repository := NewRepository(secretsManager)

		// when
		encodedCrt, encodedKey, err := repository.Get(reName)

		// then
		require.Error(t, err)
		assert.Equal(t, apperrors.CodeNotFound, err.Code())
		assert.Nil(t, encodedCrt)
		assert.Nil(t, encodedKey)
	})

	t.Run("should fail if couldn't get secret", func(t *testing.T) {
		// given
		secretsManager := &mocks.Manager{}
		secretsManager.On("Get", reName, metav1.GetOptions{}).Return(nil, &k8serrors.StatusError{})

		repository := NewRepository(secretsManager)

		// when
		encodedCrt, encodedKey, err := repository.Get(reName)

		// then
		require.Error(t, err)
		assert.Equal(t, apperrors.CodeInternal, err.Code())
		assert.Nil(t, encodedCrt)
		assert.Nil(t, encodedKey)
	})
}
