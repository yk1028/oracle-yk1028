package cli_test

import (
	"strings"
	"testing"

	"github.com/yk1028/oracle-yk1028/app"
	"github.com/yk1028/oracle-yk1028/x/oracle/client/cli"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	testnet "github.com/cosmos/cosmos-sdk/testutil/network"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     testnet.Config
	network *testnet.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := app.DefaultConfig()
	cfg.NumValidators = 1

	s.cfg = cfg
	s.network = testnet.New(s.T(), cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestGetQueryCmd() {
	val := s.network.Validators[0]

	testCases := map[string]struct {
		cmd            *cobra.Command
		args           []string
		expectedOutput string
		expectErr      bool
	}{
		"non-existent claim": {
			cli.CmdClaim(),
			[]string{"DF0C23E8634E480F84B9D5674A7CDC9816466DEC28A3358F73260F68D28D7660"},
			"claim DF0C23E8634E480F84B9D5674A7CDC9816466DEC28A3358F73260F68D28D7660 not found",
			true,
		},
		"list-claim (default pagination)": {
			cli.CmdAllClaims(),
			[]string{},
			"claims: []\npagination:\n  next_key: null\n  total: \"0\"",
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, tc.cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}

			s.Require().Contains(strings.TrimSpace(out.String()), tc.expectedOutput)
		})
	}
}
