package cmd

import (
	"net"

	"github.com/LINBIT/linstor-iscsi/pkg/crmcontrol"
	"github.com/LINBIT/linstor-iscsi/pkg/linstorcontrol"
	"github.com/LINBIT/linstor-iscsi/pkg/nfs"
	"github.com/LINBIT/linstor-iscsi/pkg/nfsbase"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// deleteCommand represents the delete command
func deleteNfsCommand() *cobra.Command {
	var controller net.IP
	var resourceName string

	var deleteCmd = &cobra.Command{
		Use:   "delete-nfs",
		Short: "Deletes an NFS export",
		Long: `Deletes an NFS export by stopping and deleting the pacemaker resource
primitives and removing the linstor resources.

For example:
linstor-iscsi delete-nfs --resource=example`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if !cmd.Flags().Changed("controller") {
				foundIP, err := crmcontrol.FindLinstorController()
				if err == nil { // it might be ok to not find it...
					controller = foundIP
				}
			}
			linstorCfg := linstorcontrol.Linstor{
				ResourceName: resourceName,
				Loglevel:     log.GetLevel().String(),
				ControllerIP: controller,
			}
			nfsCfg := nfsbase.NfsConfig{
				ResourceName: resourceName,
			}
			nfsRsc := nfs.NfsResource{
				Nfs:     nfsCfg,
				Linstor: linstorCfg,
			}
			err := nfsRsc.DeleteResource()
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	deleteCmd.Flags().StringVarP(&resourceName, "resource", "r", "", "Set the resource name (required)")
	deleteCmd.Flags().IPVarP(&controller, "controller", "c", net.IPv4(127, 0, 0, 1), "Set the IP of the linstor controller node")

	deleteCmd.MarkFlagRequired("resource")

	return deleteCmd
}