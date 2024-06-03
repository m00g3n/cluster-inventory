package azure

import (
	"github.com/gardener/gardener-extension-provider-azure/pkg/apis/azure/v1alpha1"
	"math/big"
	"net/netip"
	"strconv"
)

const workersBits = 3
const cidrLength = 32

func generateAzureZones(workerCidr string, zoneNames []string) []v1alpha1.Zone {
	var zones []v1alpha1.Zone
	defaultConnectionTimeOutMinutes := int32(4)

	cidr, _ := netip.ParsePrefix(workerCidr)
	workerPrefixLength := cidr.Bits() + workersBits
	workerPrefix, _ := cidr.Addr().Prefix(workerPrefixLength)
	// delta - it is the difference between CIDRs of two zones:
	//    zone1:   "10.250.0.0/19",
	//    zone2:   "10.250.32.0/19",
	delta := big.NewInt(1)
	delta.Lsh(delta, uint(cidrLength-workerPrefixLength))

	// zoneIPValue - it is an integer, which is based on IP bytes
	zoneIPValue := new(big.Int).SetBytes(workerPrefix.Addr().AsSlice())

	for _, name := range convertZoneNames(zoneNames) {
		zoneWorkerIP, _ := netip.AddrFromSlice(zoneIPValue.Bytes())
		zoneWorkerCidr := netip.PrefixFrom(zoneWorkerIP, workerPrefixLength)
		zoneIPValue.Add(zoneIPValue, delta)
		zones = append(zones, v1alpha1.Zone{
			Name: name,
			CIDR: zoneWorkerCidr.String(),
			NatGateway: &v1alpha1.ZonedNatGatewayConfig{
				// There are existing Azure clusters which were created before NAT gateway support,
				// and they were migrated to HA with all zones having enableNatGateway: false .
				// But for new Azure runtimes, enableNatGateway for all zones is always true
				Enabled:                      true,
				IdleConnectionTimeoutMinutes: &defaultConnectionTimeOutMinutes,
			},
		})
	}
	return zones
}

func convertZoneNames(zoneNames []string) []int32 {
	var zones []int32
	for _, inputZone := range zoneNames {
		zone, err := strconv.Atoi(inputZone)
		if err != nil || zone < 1 || zone > 3 {
			continue
		}
		zones = append(zones, int32(zone))
	}

	return zones
}
