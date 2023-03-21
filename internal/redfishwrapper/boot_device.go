package redfishwrapper

import (
	"context"

	bmclibErrs "github.com/bmc-toolbox/bmclib/v2/errors"
	"github.com/pkg/errors"
	rf "github.com/stmcginnis/gofish/redfish"
)

// Set the boot device for the system.
func (c *Client) SystemBootDeviceSet(ctx context.Context, bootDevice string, setPersistent, efiBoot bool) (ok bool, err error) {
	if err := c.SessionActive(); err != nil {
		return false, errors.Wrap(bmclibErrs.ErrNotAuthenticated, err.Error())
	}

	systems, err := c.client.Service.Systems()
	if err != nil {
		return false, err
	}

	for _, system := range systems {
		boot := system.Boot

		switch bootDevice {
		case "bios":
			boot.BootSourceOverrideTarget = rf.BiosSetupBootSourceOverrideTarget
		case "cdrom":
			boot.BootSourceOverrideTarget = rf.CdBootSourceOverrideTarget
		case "diag":
			boot.BootSourceOverrideTarget = rf.DiagsBootSourceOverrideTarget
		case "floppy":
			boot.BootSourceOverrideTarget = rf.FloppyBootSourceOverrideTarget
		case "disk":
			boot.BootSourceOverrideTarget = rf.HddBootSourceOverrideTarget
		case "none":
			boot.BootSourceOverrideTarget = rf.NoneBootSourceOverrideTarget
		case "pxe":
			boot.BootSourceOverrideTarget = rf.PxeBootSourceOverrideTarget
		case "remote_drive":
			boot.BootSourceOverrideTarget = rf.RemoteDriveBootSourceOverrideTarget
		case "sd_card":
			boot.BootSourceOverrideTarget = rf.SDCardBootSourceOverrideTarget
		case "usb":
			boot.BootSourceOverrideTarget = rf.UsbBootSourceOverrideTarget
		case "utilities":
			boot.BootSourceOverrideTarget = rf.UtilitiesBootSourceOverrideTarget
		default:
			return false, errors.New("invalid boot device")
		}

		if setPersistent {
			boot.BootSourceOverrideEnabled = rf.ContinuousBootSourceOverrideEnabled
		} else {
			boot.BootSourceOverrideEnabled = rf.OnceBootSourceOverrideEnabled
		}

		var bootMode rf.BootSourceOverrideMode
		if efiBoot {
			bootMode = rf.UEFIBootSourceOverrideMode
		} else {
			bootMode = rf.LegacyBootSourceOverrideMode
		}
		// HPE iLO treats BootSourceOverrideMode as read-only and errors even if
		// it is set to the existing value, so only attempt to set it if
		// required.
		if boot.BootSourceOverrideMode != bootMode {
			boot.BootSourceOverrideMode = bootMode
		} else {
			boot.BootSourceOverrideMode = ""
		}
		// HPE iLO will error if UefiTargetBootSourceOverride is set when
		// BootSourceOverrideTarget isn't UefiTarget. This function doesn't ever
		// set UefiTarget anyway, so clear this.
		boot.UefiTargetBootSourceOverride = ""

		err = system.SetBoot(boot)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
