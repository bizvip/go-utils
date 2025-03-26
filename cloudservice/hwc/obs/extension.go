package obs

import (
	"fmt"
	"strconv"
	"strings"
)

type extensionOptions interface{}
type extensionHeaders func(headers map[string][]string, isObs bool) error

func WithProgress(progressListener ProgressListener) configurer {
	return func(conf *config) {
		conf.progressListener = progressListener
	}
}
func setHeaderPrefix(key string, value string) extensionHeaders {
	return func(headers map[string][]string, isObs bool) error {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("set header %s with empty value", key)
		}
		setHeaders(headers, key, []string{value}, isObs)
		return nil
	}
}

// WithReqPaymentHeader sets header for requester-pays
func WithReqPaymentHeader(requester PayerType) extensionHeaders {
	return setHeaderPrefix(REQUEST_PAYER, string(requester))
}

func WithTrafficLimitHeader(trafficLimit int64) extensionHeaders {
	return setHeaderPrefix(TRAFFIC_LIMIT, strconv.FormatInt(trafficLimit, 10))
}

func WithCallbackHeader(callback string) extensionHeaders {
	return setHeaderPrefix(CALLBACK, string(callback))
}

func WithCustomHeader(key string, value string) extensionHeaders {
	return func(headers map[string][]string, isObs bool) error {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("set header %s with empty value", key)
		}
		headers[key] = []string{value}
		return nil
	}
}
