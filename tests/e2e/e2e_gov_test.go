package e2e

import (
	"context"
	"fmt"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdkgovtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (s *IntegrationTestSuite) executeSubmitUpgradeProposal(c *chain, valIdx, upgradeHeight int, submitterAddr, proposalName, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu tx submit proposal %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		sdkgovtypes.ModuleName,
		"submit-proposal",
		"software-upgrade",
		proposalName,
		fmt.Sprintf("--upgrade-height=%d", upgradeHeight),
		fmt.Sprintf("--title=\"title of %s\"", proposalName),
		fmt.Sprintf("--description=\"description of %s\"", proposalName),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, submitterAddr),
		fmt.Sprintf("--%s=%s", flags.FlagGasPrices, fees),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}

	s.T().Logf("cmd: %s", strings.Join(command, " "))

	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("%s successfully submit %s proposal", submitterAddr, proposalName)
}

func (s *IntegrationTestSuite) executeDepositProposal(c *chain, valIdx int, submitterAddr string, proposalId int, amount, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu tx deposit proposal %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		sdkgovtypes.ModuleName,
		"deposit",
		fmt.Sprintf("%d", proposalId),
		amount,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, submitterAddr),
		fmt.Sprintf("--%s=%s", flags.FlagGasPrices, fees),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}

	s.T().Logf("cmd: %s", strings.Join(command, " "))

	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("%s successfully deposit proposal %d %s", submitterAddr, proposalId, amount)
}

func (s *IntegrationTestSuite) executeVoteProposal(c *chain, valIdx int, submitterAddr string, proposalId int, vote, fees string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("Executing shentu tx vote proposal %s", c.id)

	command := []string{
		shentuBinary,
		txCommand,
		sdkgovtypes.ModuleName,
		"vote",
		fmt.Sprintf("%d", proposalId),
		vote,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, submitterAddr),
		fmt.Sprintf("--%s=%s", flags.FlagGasPrices, fees),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}

	s.T().Logf("cmd: %s", strings.Join(command, " "))

	s.execShentuTxCmd(ctx, c, command, valIdx, s.defaultExecValidation(c, valIdx))
	s.T().Logf("%s successfully vote proposal %d %s", submitterAddr, proposalId, vote)
}

func queryProposal(endpoint string, proposalID int) (*govtypesv1.QueryProposalResponse, error) {
	grpcReq := &govtypesv1.QueryProposalRequest{
		ProposalId: uint64(proposalID),
	}
	conn, err := connectGrpc(endpoint)
	defer conn.Close()
	client := govtypesv1.NewQueryClient(conn)

	grpcRsp, err := client.Proposal(context.Background(), grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return grpcRsp, nil
}
