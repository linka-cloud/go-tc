package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

func extractTcmsgAttributes(data []byte, info *Attribute) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var options []byte
	var xStats []byte
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaKind:
			info.Kind = ad.String()
		case tcaOptions:
			// the evaluation of this field depends on tcaKind.
			// there is no guarantee, that kind is know at this moment,
			// so we save it for later
			options = ad.Bytes()
		case tcaChain:
			info.Chain = ad.Uint32()
		case tcaXstats:
			// the evaluation of this field depends on tcaKind.
			// there is no guarantee, that kind is know at this moment,
			// so we save it for later
			xStats = ad.Bytes()
		case tcaStats:
			tcstats := &Stats{}
			if err := extractTCStats(ad.Bytes(), tcstats); err != nil {
				return err
			}
			info.Stats = tcstats
		case tcaStats2:
			tcstats2 := &Stats2{}
			if err := extractTCStats2(ad.Bytes(), tcstats2); err != nil {
				return err
			}
			info.Stats2 = tcstats2
		case tcaHwOffload:
			info.HwOffload = ad.Uint8()
		case tcaEgressBlock:
			info.EgressBlock = ad.Uint32()
		case tcaIngressBlock:
			info.IngressBlock = ad.Uint32()
		default:
			return fmt.Errorf("extractTcmsgAttributes()\t%d\n\t%v", ad.Type(), ad.Bytes())

		}
	}
	if len(options) > 0 {
		if err := extractTCAOptions(options, info, info.Kind); err != nil {
			return err
		}
	}
	if len(xStats) > 0 {
		tcxstats := &XStats{}
		if err := extractXStats(xStats, tcxstats, info.Kind); err != nil {
			return err
		}
		info.XStats = tcxstats
	}
	return nil
}

func extractTCAOptions(data []byte, tc *Attribute, kind string) error {
	switch kind {
	case "fq_codel":
		info := &FqCodel{}
		if err := unmarshalFqCodel(data, info); err != nil {
			return err
		}
		tc.FqCodel = info
	case "codel":
		info := &Codel{}
		if err := unmarshalCodel(data, info); err != nil {
			return err
		}
		tc.Codel = info
	case "fq":
		info := &Fq{}
		if err := unmarshalFq(data, info); err != nil {
			return err
		}
		tc.Fq = info
	case "pie":
		info := &Pie{}
		if err := unmarshalPie(data, info); err != nil {
			return err
		}
		tc.Pie = info
	case "hhf":
		info := &Hhf{}
		if err := unmarshalHhf(data, info); err != nil {
			return err
		}
		tc.Hhf = info
	case "htb":
		info := &Htb{}
		if err := unmarshalHtb(data, info); err != nil {
			return err
		}
		tc.Htb = info
	case "hfsc":
		info := &Hfsc{}
		if err := unmarshalHfsc(data, info); err != nil {
			return err
		}
		tc.Hfsc = info
	case "dsmark":
		info := &Dsmark{}
		if err := unmarshalDsmark(data, info); err != nil {
			return err
		}
		tc.Dsmark = info
	case "drr":
		info := &Drr{}
		if err := unmarshalDrr(data, info); err != nil {
			return err
		}
		tc.Drr = info
	case "cbq":
		info := &Cbq{}
		if err := unmarshalCbq(data, info); err != nil {
			return err
		}
		tc.Cbq = info
	case "atm":
		info := &Atm{}
		if err := unmarshalAtm(data, info); err != nil {
			return err
		}
		tc.Atm = info
	case "tbf":
		info := &Tbf{}
		if err := unmarshalTbf(data, info); err != nil {
			return err
		}
		tc.Tbf = info
	case "sfb":
		info := &Sfb{}
		if err := unmarshalSfb(data, info); err != nil {
			return err
		}
		tc.Sfb = info
	case "red":
		info := &Red{}
		if err := unmarshalRed(data, info); err != nil {
			return err
		}
		tc.Red = info
	case "pfifo":
		limit := &FifoOpt{}
		if err := extractFifoOpt(data, limit); err != nil {
			return err
		}
		tc.Pfifo = limit
	case "mqprio":
		info := &MqPrio{}
		if err := unmarshalMqPrio(data, info); err != nil {
			return err
		}
		tc.MqPrio = info
	case "bfifo":
		limit := &FifoOpt{}
		if err := extractFifoOpt(data, limit); err != nil {
			return err
		}
		tc.Bfifo = limit
	case "clsact":
		return extractClsact(data)
	case "ingress":
		return extractIngress(data)
	case "qfq":
		info := &Qfq{}
		if err := unmarshalQfq(data, info); err != nil {
			return err
		}
		tc.Qfq = info
	case "bpf":
		info := &Bpf{}
		if err := unmarshalBpf(data, info); err != nil {
			return err
		}
		tc.BPF = info
	case "u32":
		info := &U32{}
		if err := unmarshalU32(data, info); err != nil {
			return err
		}
		tc.U32 = info
	case "rsvp":
		info := &Rsvp{}
		if err := unmarshalRsvp(data, info); err != nil {
			return err
		}
		tc.Rsvp = info
	case "route":
		info := &Route4{}
		if err := unmarshalRoute4(data, info); err != nil {
			return err
		}
		tc.Route4 = info
	case "fw":
		info := &Fw{}
		if err := unmarshalFw(data, info); err != nil {
			return err
		}
		tc.Fw = info
	case "flow":
		info := &Flow{}
		if err := unmarshalFlow(data, info); err != nil {
			return err
		}
		tc.Flow = info
	default:
		return fmt.Errorf("extractTCAOptions(): unsupported kind: %s", kind)
	}

	return nil
}

func extractXStats(data []byte, tc *XStats, kind string) error {
	switch kind {
	case "sfq":
		info := &SfqXStats{}
		if err := extractSfqXStats(data, info); err != nil {
			return err
		}
		tc.Sfq = info
	case "sfb":
		info := &SfbXStats{}
		if err := extractSfbXStats(data, info); err != nil {
			return err
		}
		tc.Sfb = info
	case "red":
		info := &RedXStats{}
		if err := extractRedXStats(data, info); err != nil {
			return err
		}
		tc.Red = info
	case "choke":
		info := &ChokeXStats{}
		if err := extractChokeXStats(data, info); err != nil {
			return err
		}
		tc.Choke = info
	case "htb":
		info := &HtbXStats{}
		if err := extractHtbXStats(data, info); err != nil {
			return err
		}
		tc.Htb = info
	case "cbq":
		info := &CbqXStats{}
		if err := extractCbqXStats(data, info); err != nil {
			return err
		}
		tc.Cbq = info
	case "codel":
		info := &CodelXStats{}
		if err := extractCodelXStats(data, info); err != nil {
			return err
		}
		tc.Codel = info
	case "hhf":
		info := &HhfXStats{}
		if err := extractHhfXStats(data, info); err != nil {
			return err
		}
		tc.Hhf = info
	case "pie":
		info := &PieXStats{}
		if err := extractPieXStats(data, info); err != nil {
			return err
		}
		tc.Pie = info
	case "fq_codel":
		info := &FqCodelXStats{}
		if err := extractFqCodelXStats(data, info); err != nil {
			return err
		}
		tc.FqCodel = info
	default:
		return fmt.Errorf("extractXStats(): unsupported kind: %s", kind)
	}
	return nil
}

func extractClsact(data []byte) error {
	// Clsact is parameterless - so we expect to options
	if len(data) != 0 {
		return fmt.Errorf("extractClsact()\t%v", data)
	}
	return nil
}

func extractIngress(data []byte) error {
	// Ingress is parameterless - so we expect to options
	if len(data) != 0 {
		return fmt.Errorf("extractIngress()\t%v", data)
	}
	return nil
}

const (
	tcaUnspec = iota
	tcaKind
	tcaOptions
	tcaStats
	tcaXstats
	tcaRate
	tcaFcnt
	tcaStats2
	tcaStab
	tcaPad
	tcaDumpInvisible
	tcaChain
	tcaHwOffload
	tcaIngressBlock
	tcaEgressBlock
)

const (
	tcaStatsUnspec = iota
	tcaStatsBasic
	tcaStatsRateEst
	tcaStatsQueue
	tcaStatsApp
	tcaStatsRateEst64
	tcaStatsPAD
	tcaStatsBasicHw
)
