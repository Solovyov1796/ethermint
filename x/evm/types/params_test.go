package types

import (
	"testing"

	"github.com/ethereum/go-ethereum/params"

	"github.com/stretchr/testify/require"
)

func TestParamsValidate(t *testing.T) {
	extraEips := []int64{2929, 1884, 1344}
	testCases := []struct {
		name     string
		params   Params
		expError bool
	}{
		{"default", DefaultParams(), false},
		{
			"valid",
			NewParams("ara", false, true, true, DefaultChainConfig(), extraEips, []EIP712AllowedMsg{}),
			false,
		},
		{
			"empty",
			Params{},
			true,
		},
		{
			"invalid evm denom",
			Params{
				EvmDenom: "@!#!@$!@5^32",
			},
			true,
		},
		{
			"invalid eip",
			Params{
				EvmDenom:  "stake",
				ExtraEIPs: []int64{1},
			},
			expError: true,
		},
		{
			name: "invalid eip",
			getParams: func() Params {
				return Params{
					EvmDenom:  "stake",
					ExtraEIPs: []int64{1},
				}
			},
			expError: true,
		},
	}

	for _, tc := range testCases {
		err := tc.params.Validate()

		if tc.expError {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}

func TestEnabledPrecompilesAddressCorrectness(t *testing.T) {
	const (
		validEthAddress   = "0xc0ffee254729296a45a3885639AC7E10F9d54979"
		invalidEthAddress = "0xc0ffee254729296a45a3885639AC7E10F9d5497"
	)

	testCases := []struct {
		name      string
		getParams func() Params
		errorMsg  string
	}{
		{
			name: "failure: empty address",
			getParams: func() Params {
				params := DefaultParams()
				params.EnabledPrecompiles = []string{""}
				return params
			},
			errorMsg: "invalid hex address",
		},
		{
			name: "failure: invalid address #1",
			getParams: func() Params {
				params := DefaultParams()
				params.EnabledPrecompiles = []string{invalidEthAddress}
				return params
			},
			errorMsg: "invalid hex address",
		},
		{
			name: "failure: invalid address #2",
			getParams: func() Params {
				params := DefaultParams()
				params.EnabledPrecompiles = []string{validEthAddress, invalidEthAddress}
				return params
			},
			errorMsg: "invalid hex address",
		},
		{
			name: "success: pass nil as enabled precompiles",
			getParams: func() Params {
				params := DefaultParams()
				params.EnabledPrecompiles = nil
				return params
			},
			errorMsg: "",
		},
		{
			name: "success: pass empty slice as enabled precompiles",
			getParams: func() Params {
				params := DefaultParams()
				params.EnabledPrecompiles = []string{}
				return params
			},
			errorMsg: "",
		},
		{
			name: "success: valid address",
			getParams: func() Params {
				params := DefaultParams()
				params.EnabledPrecompiles = []string{validEthAddress}
				return params
			},
			errorMsg: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.getParams().Validate()

			if tc.errorMsg == "" {
				require.NoError(t, err, tc.name)
			} else {
				require.ErrorContains(t, err, tc.errorMsg, tc.name)
			}
		})
	}
}

func TestEnabledPrecompilesOrderInBytesRepr(t *testing.T) {
	const (
		addr1 = "0x1000000000000000000000000000000000000000"
		addr2 = "0x2000000000000000000000000000000000000000"

		// NOTE: we sort in bytes representation, so proper order will be []string{mixedCaseAddr, upperCaseAddr},
		// and it differs from lexicographically sorted strings
		upperCaseAddr = "0xAB00000000000000000000000000000000000000"
		mixedCaseAddr = "0xaA00000000000000000000000000000000000000"
	)

	testCases := []struct {
		name      string
		getParams func() Params
		errorMsg  string
	}{
		{
			name: "success: addresses are sorted",
			getParams: func() Params {
				params := DefaultParams()
				params.EnabledPrecompiles = []string{addr1, addr2}
				return params
			},
			errorMsg: "",
		},
		{
			name: "failure: addresses are in reverse order",
			getParams: func() Params {
				params := DefaultParams()
				params.EnabledPrecompiles = []string{addr2, addr1}
				return params
			},
			errorMsg: "enabled precompiles are not sorted",
		},
		{
			name: "success: addresses are sorted in bytes representation",
			getParams: func() Params {
				params := DefaultParams()
				params.EnabledPrecompiles = []string{mixedCaseAddr, upperCaseAddr}
				return params
			},
			errorMsg: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.getParams().Validate()

			if tc.errorMsg == "" {
				require.NoError(t, err, tc.name)
			} else {
				require.ErrorContains(t, err, tc.errorMsg, tc.name)
			}
		})
	}
}

func TestEnabledPrecompilesUniquenessInBytesRepr(t *testing.T) {
	const (
		addr1 = "0x1000000000000000000000000000000000000000"
		addr2 = "0x2000000000000000000000000000000000000000"

		// NOTE: we check uniqueness in bytes representation, so lowerCaseAddr and mixedCaseAddr are the same,
		// despite it differs in string representation
		lowerCaseAddr = "0xab00000000000000000000000000000000000000"
		mixedCaseAddr = "0xAb00000000000000000000000000000000000000"
	)

	testCases := []struct {
		name      string
		getParams func() Params
		errorMsg  string
	}{
		{
			name: "success: addresses are unique",
			getParams: func() Params {
				params := DefaultParams()
				params.EnabledPrecompiles = []string{addr1, addr2}
				return params
			},
			errorMsg: "",
		},
		{
			name: "failure: addresses are not unique",
			getParams: func() Params {
				params := DefaultParams()
				params.EnabledPrecompiles = []string{addr1, addr1}
				return params
			},
			errorMsg: "enabled precompiles are not unique",
		},
		{
			name: "failure: addresses are not unique in bytes representation",
			getParams: func() Params {
				params := DefaultParams()
				params.EnabledPrecompiles = []string{lowerCaseAddr, mixedCaseAddr}
				return params
			},
			errorMsg: "enabled precompiles are not unique",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.getParams().Validate()

			if tc.errorMsg == "" {
				require.NoError(t, err, tc.name)
			} else {
				require.ErrorContains(t, err, tc.errorMsg, tc.name)
			}
		})
	}
}

func TestParamsEIPs(t *testing.T) {
	extraEips := []int64{2929, 1884, 1344}
	params := NewParams("ara", false, true, true, DefaultChainConfig(), extraEips, []EIP712AllowedMsg{})
	actual := params.EIPs()

	require.Equal(t, []int([]int{2929, 1884, 1344}), actual)
}

func TestParamsValidatePriv(t *testing.T) {
	require.Error(t, validateEVMDenom(false))
	require.NoError(t, validateEVMDenom("inj"))
	require.Error(t, validateBool(""))
	require.NoError(t, validateBool(true))
	require.Error(t, validateEIPs(""))
	require.NoError(t, validateEIPs([]int64{1884}))
}

func TestValidateChainConfig(t *testing.T) {
	testCases := []struct {
		name     string
		i        interface{}
		expError bool
	}{
		{
			"invalid chain config type",
			"string",
			true,
		},
		{
			"valid chain config type",
			DefaultChainConfig(),
			false,
		},
	}
	for _, tc := range testCases {
		err := validateChainConfig(tc.i)

		if tc.expError {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}

func TestIsLondon(t *testing.T) {
	testCases := []struct {
		name   string
		height int64
		result bool
	}{
		{
			"Before london block",
			5,
			false,
		},
		{
			"After london block",
			12_965_001,
			true,
		},
		{
			"london block",
			12_965_000,
			true,
		},
	}

	for _, tc := range testCases {
		ethConfig := params.MainnetChainConfig
		require.Equal(t, IsLondon(ethConfig, tc.height), tc.result)
	}
}
