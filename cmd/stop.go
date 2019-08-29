package cmd

import (
	"github.com/LINBIT/linstor-iscsi/pkg/iscsi"
	"github.com/LINBIT/linstor-iscsi/pkg/linstorcontrol"
	"github.com/LINBIT/linstor-iscsi/pkg/targetutil"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops an iSCSI target",
	Long: `Sets the target role attribute of a Pacemaker primitive to stopped.
This causes pacemaker to stop the components of an iSCSI target.

For example:
linstor-iscsi start --iqn=iqn.2019-08.com.linbit:example --lun=1`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		linstorCfg := linstorcontrol.Linstor{
			Loglevel:     log.GetLevel().String(),
			ControllerIP: controller,
		}
		targetCfg := targetutil.TargetConfig{
			IQN:  iqn,
			LUNs: []*targetutil.LUN{&targetutil.LUN{ID: uint8(lun)}},
		}
		target := cliNewTargetMust(cmd, targetCfg)
		iscsiCfg := &iscsi.ISCSI{Linstor: linstorCfg, Target: target}
		if err := iscsiCfg.StopResource(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	stopCmd.MarkPersistentFlagRequired("iqn")
	stopCmd.MarkPersistentFlagRequired("lun")
}
