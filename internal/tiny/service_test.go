package tiny_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"sanathk.com/tinyurl/internal/tiny"
	"sanathk.com/tinyurl/pkg/fakes"
)

func TestServiceInit(t *testing.T) {
	fakeDB := fakes.FakePostgresInterface{}

	_, err := tiny.NewService(context.Background(), &fakeDB)
	require.NoError(t, err)
}

func TestServiceCreateTinyURL(t *testing.T) {
	fakeDB := fakes.FakePostgresInterface{}

	svc, err := tiny.NewService(context.Background(), &fakeDB)
	require.NoError(t, err)

	type testcase struct {
		name     string
		in       tiny.TinyURLRequest
		expected tiny.TinyURLResponse
	}

	validExpiryTime := time.Now().Add(10 * time.Hour)
	validOriginalLink := "https://www.sanathk.com"

	testcases := []testcase{
		{
			name: "valid expiry and original",
			in: tiny.TinyURLRequest{
				Expiry:   validExpiryTime,
				Original: validOriginalLink,
			},
			expected: tiny.TinyURLResponse{
				Expiry:   validExpiryTime,
				Original: validOriginalLink,
			},
			// TODO: add more tests:
			// Negative Cases:
			// 	* invalid time
			//  * invalid or empty url
			//  * max url length for hashing?
		},
	}

	for _, tc := range testcases {
		resp, err := svc.CreateTinyURL(context.Background(), tc.in)
		require.NoError(t, err)

		require.Equal(t, tc.expected.Original, resp.Original)
		require.Equal(t, tc.expected.Expiry, resp.Expiry)
		require.NotEqual(t, tc.expected.Original, resp.Tinyurl)
	}
}
