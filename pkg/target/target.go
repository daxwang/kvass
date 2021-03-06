package target

import (
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/pkg/labels"
	"net/url"
	"strings"
)

const (
	// PrefixForInvalidLabelName is a prefix string for mark invalid label name become valid
	PrefixForInvalidLabelName = model.ReservedLabelPrefix + "invalid_label_"
)

// Target is a target generate prometheus config
type Target struct {
	// Hash is calculated from origin labels before relabel_configs process and the URL of this target
	// see prometheus scrape.Target.hash
	Hash uint64 `json:"hash"`
	// Labels is result of relabel_configs process
	Labels labels.Labels `json:"labels"`
	// Series is reference series of this target, may from target explorer
	Series int64 `json:"series"`
}

// Address return the address from labels
func (t *Target) Address() string {
	for _, v := range t.Labels {
		if v.Name == model.AddressLabel {
			return v.Value
		}
	}
	return ""
}

// NoReservedLabel return the labels without reserved prefix "__"
func (t *Target) NoReservedLabel() labels.Labels {
	lset := make(labels.Labels, 0, len(t.Labels))
	for _, l := range t.Labels {
		if !strings.HasPrefix(l.Name, model.ReservedLabelPrefix) {
			lset = append(lset, l)
		}
	}
	return lset
}

// NoParamURL return a url without params
func (t *Target) NoParamURL() *url.URL {
	return &url.URL{
		Scheme: t.Labels.Get(model.SchemeLabel),
		Host:   t.Labels.Get(model.AddressLabel),
		Path:   t.Labels.Get(model.MetricsPathLabel),
	}
}

// URL return the full url of this target, the params of cfg will be add to url
func (t *Target) URL(cfg *config.ScrapeConfig) *url.URL {
	params := url.Values{}

	for k, v := range cfg.Params {
		params[k] = make([]string, len(v))
		copy(params[k], v)
	}
	for _, l := range t.Labels {
		if !strings.HasPrefix(l.Name, model.ParamLabelPrefix) {
			continue
		}
		ks := l.Name[len(model.ParamLabelPrefix):]

		if len(params[ks]) > 0 {
			params[ks][0] = l.Value
		} else {
			params[ks] = []string{l.Value}
		}
	}

	return &url.URL{
		Scheme:   t.Labels.Get(model.SchemeLabel),
		Host:     t.Labels.Get(model.AddressLabel),
		Path:     t.Labels.Get(model.MetricsPathLabel),
		RawQuery: params.Encode(),
	}
}
