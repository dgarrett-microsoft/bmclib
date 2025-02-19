package redfish

import (
	"context"
)

func (c *Conn) GetBiosConfiguration(ctx context.Context) (biosConfig map[string]string, err error) {
	systems, err := c.redfishwrapper.Systems()
	if err != nil {
		return nil, err
	}

	biosConfig = make(map[string]string)
	for _, sys := range systems {
		if !compatibleOdataID(sys.ODataID, systemsOdataIDs) {
			continue
		}

		bios, err := sys.Bios()
		if err != nil {
			return nil, err
		}

		for attr := range bios.Attributes {
			biosConfig[attr] = bios.Attributes.String(attr)
		}
	}
	return biosConfig, nil
}
