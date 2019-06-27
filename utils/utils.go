// Copyright © 2018 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

//NopHandler is an empty handler to help net/http -> Gin conversions
type NopHandler struct{}

func (h NopHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {}

//WriteToFile write the []byte to the given file
func WriteToFile(data []byte, file string) error {
	if err := os.MkdirAll(filepath.Dir(file), os.ModePerm); err != nil {
		return err
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return ioutil.WriteFile(file, data, 0644)
	}

	tmpfi, err := ioutil.TempFile(filepath.Dir(file), "file.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfi.Name())

	if err = ioutil.WriteFile(tmpfi.Name(), data, 0644); err != nil {
		return err
	}

	if err = tmpfi.Close(); err != nil {
		return err
	}

	if err = os.Remove(file); err != nil {
		return err
	}

	err = os.Rename(tmpfi.Name(), file)
	return err
}

// ConvertJson2Map converts []byte to map[string]string
func ConvertJson2Map(js []byte) (map[string]string, error) {
	var result map[string]string
	err := json.Unmarshal(js, &result)
	return result, err
}

// Contains checks slice contains `s` string
func Contains(slice []string, s string) bool {
	for _, sl := range slice {
		if sl == s {
			return true
		}
	}
	return false
}

// EncodeStringToBase64 first checks if the string is encoded if yes returns it if no than encodes it.
func EncodeStringToBase64(s string) string {
	if _, err := base64.StdEncoding.DecodeString(s); err != nil {
		return base64.StdEncoding.EncodeToString([]byte(s))
	}
	return s
}


// AddressCount returns the number of distinct host addresses within the given
// CIDR range.
func AddressCount(network *net.IPNet) uint64 {
	prefixLen, bits := network.Mask.Size()
	return 1 << (uint64(bits) - uint64(prefixLen))
}

//VerifyNoOverlap takes a list subnets and supernet (CIDRBlock) and verifies
//none of the subnets overlap and all subnets are in the supernet
//it returns an error if any of those conditions are not satisfied
func VerifyNoOverlapWithinSupernet(subnets []*net.IPNet, CIDRBlock *net.IPNet) error {
	firstLastIP := make([][]net.IP, len(subnets))
	for i, s := range subnets {
		first, last := AddressRange(s)
		firstLastIP[i] = []net.IP{first, last}
	}
	for i, s := range subnets {
		if !CIDRBlock.Contains(firstLastIP[i][0]) || !CIDRBlock.Contains(firstLastIP[i][1]) {
			return fmt.Errorf("%s does not fully contain %s", CIDRBlock.String(), s.String())
		}
		for j := 0; j < len(subnets); j++ {
			if i == j {
				continue
			}

			first := firstLastIP[j][0]
			last := firstLastIP[j][1]
			if s.Contains(first) || s.Contains(last) {
				return fmt.Errorf("%s overlaps with %s", subnets[j].String(), s.String())
			}
		}
	}
	return nil
}

//VerifyNoOverlap takes a list subnets and verifies none of the subnets overlap
//it returns an error if any of those conditions are not satisfied
func VerifyNoOverlap(subnets []*net.IPNet) error {
	firstLastIP := make([][]net.IP, len(subnets))
	for i, s := range subnets {
		first, last := AddressRange(s)
		firstLastIP[i] = []net.IP{first, last}
	}
	for i, s := range subnets {
		for j := 0; j < len(subnets); j++ {
			if i == j {
				continue
			}

			first := firstLastIP[j][0]
			last := firstLastIP[j][1]
			if s.Contains(first) || s.Contains(last) {
				return fmt.Errorf("%s overlaps with %s", subnets[j].String(), s.String())
			}
		}
	}
	return nil
}

// AddressRange returns the first and last addresses in the given CIDR range.
func AddressRange(network *net.IPNet) (net.IP, net.IP) {
	firstIP := network.IP

	// the last IP is the network address OR NOT the mask address
	prefixLen, bits := network.Mask.Size()
	if prefixLen == bits {
		// make sure that our two slices are distinct, since they
		// would be in all other cases.
		lastIP := make([]byte, len(firstIP))
		copy(lastIP, firstIP)
		return firstIP, lastIP
	}

	firstIPInt, bits := ipToInt(firstIP)
	hostLen := uint(bits) - uint(prefixLen)
	lastIPInt := big.NewInt(1)
	lastIPInt.Lsh(lastIPInt, hostLen)
	lastIPInt.Sub(lastIPInt, big.NewInt(1))
	lastIPInt.Or(lastIPInt, firstIPInt)

	return firstIP, intToIP(lastIPInt, bits)
}


func ipToInt(ip net.IP) (*big.Int, int) {
	val := &big.Int{}
	val.SetBytes(ip)
	if len(ip) == net.IPv4len {
		return val, 32
	} else if len(ip) == net.IPv6len {
		return val, 128
	} else {
		panic(fmt.Errorf("Unsupported address length %d", len(ip)))
	}
}

func intToIP(ipInt *big.Int, bits int) net.IP {
	ipBytes := ipInt.Bytes()
	ret := make([]byte, bits/8)
	// Pack our IP bytes into the end of the return array,
	// since big.Int.Bytes() removes front zero padding.
	for i := 1; i <= len(ipBytes); i++ {
		ret[len(ret)-i] = ipBytes[len(ipBytes)-i]
	}
	return ret
}
