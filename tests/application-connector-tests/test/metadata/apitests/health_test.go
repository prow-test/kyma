package apitests

import (
	"github.com/kyma-project/kyma/tests/application-connector-tests/test/metadata/testkit"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestMetadataHealth(t *testing.T) {

	config, err := testkit.ReadConfig()
	require.NoError(t, err)

	t.Run("Application Connector Metadata", func(t *testing.T) {

		t.Run("should be healthy", func(t *testing.T) {
			// given
			url := config.MetadataServiceUrl + "/v1/health"

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
