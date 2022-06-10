package helpers

import (
	"net"
	"time"
)

func MxLookup(domain string) ([]*net.MX, error) {
	mxrecords, err := net.LookupMX(domain)

	return mxrecords, Retry(15, time.Second, func() error {
		if err != nil {
			return err
		}
		return nil
	})
}
