package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaCgroupUnspec = iota
	tcaCgroupAct
	tcaCgroupPolice
	tcaCgroupEmatches
)

// Cgroup contains attributes of the cgroup discipline
type Cgroup struct {
	Action *Action
}

// marshalCgroup returns the binary encoding of Cgroup
func marshalCgroup(info *Cgroup) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Cgroup: %w", ErrNoArg)
	}

	// TODO: improve logic and check combinations
	if info.Action != nil {
		data, err := marshalAction(info.Action)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaCgroupAct, Data: data})

	}
	return marshalAttributes(options)
}

// unmarshalCgroup parses the Cgroup-encoded data and stores the result in the value pointed to by info.
func unmarshalCgroup(data []byte, info *Cgroup) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaCgroupAct:
			act := &Action{}
			if err := unmarshalAction(ad.Bytes(), act); err != nil {
				return err
			}
			info.Action = act
		default:
			return fmt.Errorf("unmarshalCgroup()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}