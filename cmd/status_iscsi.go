package cmd

import (
	"fmt"
	"net"

	"github.com/LINBIT/linstor-gateway/pkg/crmcontrol"
	"github.com/LINBIT/linstor-gateway/pkg/iscsi"
	"github.com/LINBIT/linstor-gateway/pkg/linstorcontrol"
	"github.com/LINBIT/linstor-gateway/pkg/targetutil"
	"github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func statusISCSICommand() *cobra.Command {
	var controller net.IP
	var iqn string
	var lun int

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Displays the status of an iSCSI target or logical unit",
		Long: `Triggers Pacemaker to probe the resoruce primitives of this iSCSI target.
That means Pacemaker will run the status operation on the nodes where the
resource can run.
This makes sure that Pacemakers view of the world is updated to the state
of the world.

For example:
linstor-iscsi status --iqn=iqn.2019-08.com.linbit:example --lun=1`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if !cmd.Flags().Changed("controller") {
				foundIP, err := crmcontrol.FindLinstorController()
				if err == nil { // it might be ok to not find it...
					controller = foundIP
				}
			}
			targetCfg := targetutil.TargetConfig{
				IQN:  iqn,
				LUNs: []*targetutil.LUN{&targetutil.LUN{ID: uint8(lun)}},
			}
			target := cliNewTargetMust(cmd, targetCfg)
			targetName, err := targetutil.ExtractTargetName(targetCfg.IQN)
			if err != nil {
				log.Fatal(err)
			}
			linstorCfg := linstorcontrol.Linstor{
				Loglevel:     log.GetLevel().String(),
				ControllerIP: controller,
				ResourceName: linstorcontrol.ResourceNameFromLUN(targetName, uint8(lun)),
			}
			iscsiCfg := &iscsi.ISCSI{Linstor: linstorCfg, Target: target}

			state, err := iscsiCfg.ProbeResource()
			if err != nil {
				log.Fatal(err)
			}

			linstorState, err := linstorCfg.AggregateResourceState()
			if err != nil {
				log.Warning(err)
				linstorState = linstorcontrol.Unknown
			}

			fmt.Printf("Status of target %s, LUN %d:\n", aurora.Cyan(iqn),
				aurora.Cyan(lun))
			fmt.Printf("- Target: %s\n", stateToLongStatus(state.TargetState))
			fmt.Printf("- LU: %s\n", stateToLongStatus(state.LUStates[uint8(lun)]))
			fmt.Printf("- IP: %s\n", stateToLongStatus(state.IPState))
			fmt.Printf("- LINSTOR: %s\n", linstorStateToLongStatus(linstorState))
		},
	}

	statusCmd.Flags().StringVarP(&iqn, "iqn", "i", "", "Set the iSCSI Qualified Name (e.g., iqn.2019-08.com.linbit:unique) (required)")
	statusCmd.Flags().IntVarP(&lun, "lun", "l", 1, "Set the iSCSI LU number (LUN)")
	statusCmd.Flags().IPVarP(&controller, "controller", "c", net.IPv4(127, 0, 0, 1), "Set the IP of the linstor controller node")
	statusCmd.MarkFlagRequired("iqn")
	statusCmd.MarkFlagRequired("lun")

	return statusCmd
}
