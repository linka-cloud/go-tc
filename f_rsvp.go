package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaRsvpUnspec = iota
	tcaRsvpClassID
	tcaRsvpDst
	tcaRsvpSrc
	tcaRsvpPInfo
	tcaRsvpPolice
	tcaRsvpAct
)

// Rsvp contains attributes of the rsvp discipline
type Rsvp struct {
	ClassID *uint32
	Dst     *[]byte
	Src     *[]byte
	PInfo   *RsvpPInfo
	Police  *Police
}

// unmarshalRsvp parses the Rsvp-encoded data and stores the result in the value pointed to by info.
func unmarshalRsvp(data []byte, info *Rsvp) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaRsvpClassID:
			info.ClassID = uint32Ptr(ad.Uint32())
		case tcaRsvpDst:
			info.Dst = bytesPtr(ad.Bytes())
		case tcaRsvpSrc:
			info.Src = bytesPtr(ad.Bytes())
		case tcaRsvpPInfo:
			arg := &RsvpPInfo{}
			err := unmarshalStruct(ad.Bytes(), arg)
			concatError(multiError, err)
			info.PInfo = arg
		case tcaRsvpPolice:
			pol := &Police{}
			err := unmarshalPolice(ad.Bytes(), pol)
			concatError(multiError, err)
			info.Police = pol
		default:
			return fmt.Errorf("unmarshalRsvp()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalRsvp returns the binary encoding of Rsvp
func marshalRsvp(info *Rsvp) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Ipt: %w", ErrNoArg)
	}
	var multiError error

	// TODO: improve logic and check combinations
	if info.ClassID != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaRoute4ClassID, Data: uint32Value(info.ClassID)})
	}
	if info.PInfo != nil {
		data, err := marshalStruct(info.PInfo)
		concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaRsvpPInfo, Data: data})
	}
	if info.Src != nil {
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaRsvpSrc, Data: bytesValue(info.Src)})
	}
	if info.Dst != nil {
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaRsvpDst, Data: bytesValue(info.Dst)})
	}
	if info.Police != nil {
		data, err := marshalPolice(info.Police)
		concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaRsvpPolice, Data: data})
	}

	if multiError != nil {
		return []byte{}, multiError
	}

	return marshalAttributes(options)
}

// RsvpPInfo from include/uapi/linux/pkt_sched.h
type RsvpPInfo struct {
	Dpi       RsvpGpi
	Spi       RsvpGpi
	Protocol  uint8
	TunnelID  uint8
	TunnelHdr uint8
	Pad       uint8
}

// RsvpGpi from include/uapi/linux/pkt_sched.h
type RsvpGpi struct {
	Key    uint32
	Mask   uint32
	Offset uint32
}
