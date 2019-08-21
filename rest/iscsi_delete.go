package rest

import (
	"net/http"
	"strconv"

	"github.com/LINBIT/linstor-remote-storage/iscsi"
	"github.com/LINBIT/linstor-remote-storage/linstorcontrol"
	"github.com/gorilla/mux"
)

// ISCSIDelete deletes a highly-available iSCSI target via the REST-API
func ISCSIDelete(w http.ResponseWriter, r *http.Request) {
	tgt := mux.Vars(r)["target"]
	lid, err := strconv.Atoi(mux.Vars(r)["lun"])
	if err != nil {
		_, _ = Errorf(http.StatusBadRequest, w, "Could not convert LUN to number: %v", err)
		return
	}

	lun := iscsi.LUN{ID: uint8(lid)}
	targetConfig := iscsi.TargetConfig{
		IQN:  "iqn.1981-09.at.rck:" + tgt,
		LUNs: []*iscsi.LUN{&lun},
	}
	target, err := iscsi.NewTarget(targetConfig)
	if err != nil {
		_, _ = Errorf(http.StatusInternalServerError, w, "Could not create target from target config: %v", err)
		return
	}

	iscsiCfg := iscsi.ISCSI{
		Target:  target,
		Linstor: linstorcontrol.Linstor{},
	}

	maybeSetLinstorController(&iscsiCfg)

	if err := iscsi.CheckIQN(iscsiCfg.Target.IQN); err != nil {
		_, _ = Errorf(http.StatusBadRequest, w, "Could not validate IQN: %v", err)
		return
	}
	if err := iscsiCfg.DeleteResource(); err != nil {
		_, _ = Errorf(http.StatusInternalServerError, w, "Could not delete resource: %v", err)
		return
	}

	// json.NewEncoder(w).Encode(xxx)
}
