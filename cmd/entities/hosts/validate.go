package hosts

import "fmt"

var isValidMask = func(m int) bool { return m >= 26 && m <= 29 }

func validateNetworkArgs(networkType, distributionMethod string, mask int) error {
	if networkType != "public" && networkType != "private" {
		return fmt.Errorf("--type must be 'public' or 'private'")
	}

	if networkType == "private" {
		if distributionMethod != "gateway" {
			return fmt.Errorf("--distribution-method for private network can only be 'gateway'")
		}
		if !isValidMask(mask) {
			return fmt.Errorf("--mask for private network must be: 26, 27, 28, or 29")
		}
		return nil
	}

	if networkType == "public" {
		switch distributionMethod {
		case "gateway":
			if !isValidMask(mask) {
				return fmt.Errorf("--mask for public network (gateway) must be: 26, 27, 28, or 29")
			}
			return nil
		case "route":
			if mask != 32 {
				return fmt.Errorf("--mask for public network with distribution-method='route' must be 32")
			}
			return nil
		default:
			return fmt.Errorf("--distribution-method for public network must be 'gateway' or 'route'")
		}
	}

	return nil
}
