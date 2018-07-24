package apitests

import (
	"github.com/kyma-project/kyma/tests/gateway-tests/test/testkit"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestGatewayHealth(t *testing.T) {

	config, err := testkit.ReadConfig()
	require.NoError(t, err)

	t.Run("SF Gateway", func(t *testing.T) {

		t.Run("should be healthy", func(t *testing.T) {
			// given
			url := config.GatewayUrl + "/v1/health"

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// when
			response, err := http.DefaultClient.Do(request)
			require.NoError(t, err)

			// then
			require.Equal(t, response.StatusCode, http.StatusOK)
		})

	})

}
