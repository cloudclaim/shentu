package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	bountyQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the bounty module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	bountyQueryCmd.AddCommand(
		GetCmdQueryProgram(),
		GetCmdQueryPrograms(),
		GetCmdQueryFinding(),
		GetCmdQueryFindings(),
		GetCmdQueryFindingFingerprint(),
		GetCmdQueryProgramFingerprint(),
	)

	return bountyQueryCmd
}

// GetCmdQueryProgram implements the query program command.
func GetCmdQueryProgram() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "program [program-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of a single program",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a program. You can find the program-id by running "%s query bounty program".
Example:
$ %s query bounty program 1
`,
				version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// Query the program
			res, err := queryClient.Program(
				cmd.Context(),
				&types.QueryProgramRequest{ProgramId: args[0]},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Program)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryPrograms implements the query programs command.
func GetCmdQueryPrograms() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "programs",
		Short: "Query all programs",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for a all paginated programs that match optional filters.

Example:
$ %s query bounty programs --page=1 --limit=100
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Programs(
				cmd.Context(),
				&types.QueryProgramsRequest{
					Pagination: pageReq,
				})
			if err != nil {
				return err
			}

			if len(res.GetPrograms()) == 0 {
				return fmt.Errorf("no programs found")
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "programs")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryFinding implements the query finding command.
func GetCmdQueryFinding() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "finding [finding-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of a single finding",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a finding. You can find the finding-id by running "%s query bounty findings".
Example:
$ %s query bounty finding 1
`,
				version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// Query the finding
			res, err := queryClient.Finding(
				cmd.Context(),
				&types.QueryFindingRequest{
					FindingId: args[0],
				})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Finding)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryFindings implements the query findings command.
func GetCmdQueryFindings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "findings",
		Short: "Query findings with optional filters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for a all paginated findings that match optional filters.

Example:
$ %s query bounty findings --program-id 1
$ %s query bounty findings --submitter-address cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
$ %s query bounty findings --page=1 --limit=100
`,
				version.AppName, version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			// validate that the program-id
			pid, err := cmd.Flags().GetString(FlagProgramID)
			if err != nil {
				return err
			}

			submitterAddr, err := cmd.Flags().GetString(FlagSubmitterAddress)
			if err != nil {
				return err
			}
			if len(submitterAddr) != 0 {
				_ = sdk.MustAccAddressFromBech32(submitterAddr)
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryFindingsRequest{
				SubmitterAddress: submitterAddr,
				Pagination:       pageReq,
			}
			if len(pid) != 0 {
				req.ProgramId = pid
			}

			if len(req.ProgramId) == 0 && len(req.SubmitterAddress) == 0 {
				return fmt.Errorf("invalid request")
			}
			res, err := queryClient.Findings(cmd.Context(), req)
			if err != nil {
				return err
			}

			if len(res.GetFindings()) == 0 {
				return fmt.Errorf("no finding found")
			}

			return clientCtx.PrintProto(res)

		},
	}

	cmd.Flags().String(FlagProgramID, "", "(optional) filter by programs find by program id")
	cmd.Flags().String(FlagSubmitterAddress, "", "(optional) filter by programs find by submitter address")
	flags.AddPaginationFlagsToCmd(cmd, "findings")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryFindingFingerprint implements the query finding fingerPrint command.
func GetCmdQueryFindingFingerprint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "finding-fingerprint [finding-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query fingerPrint of a single finding",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query fingerPrint for a finding. You can find the finding-id by running "%s query bounty findings".
Example:
$ %s query bounty finding 1
`,
				version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// Query the finding
			res, err := queryClient.FindingFingerprint(
				cmd.Context(),
				&types.QueryFindingFingerprintRequest{
					FindingId: args[0],
				})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryProgramFingerprint implements the query program fingerPrint command.
func GetCmdQueryProgramFingerprint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "program-fingerprint [program-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query fingerPrint of a single program",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query fingerPrint for a program. You can find the program-id by running "%s query bounty findings".
Example:
$ %s query bounty program 1
`,
				version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// Query the program
			res, err := queryClient.ProgramFingerprint(
				cmd.Context(),
				&types.QueryProgramFingerprintRequest{
					ProgramId: args[0],
				})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
