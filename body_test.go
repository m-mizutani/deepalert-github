package main_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m-mizutani/deepalert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	main "github.com/m-mizutani/deepalert-github"
)

func TestBodyBuild(t *testing.T) {
	reportID := deepalert.ReportID(uuid.New().String())
	report := deepalert.Report{
		Result: deepalert.ReportResult{
			Severity: deepalert.SevUnclassified,
			Reason:   "It's test",
		},
		ID: reportID,
		Alerts: []deepalert.Alert{
			{
				Detector:    "blue",
				RuleName:    "orange",
				AlertKey:    "five",
				Description: "not sane",
				Timestamp:   time.Now(),
				Attributes: []deepalert.Attribute{
					{
						Type:    deepalert.TypeIPAddr,
						Key:     "source",
						Value:   "192.168.0.1",
						Context: []deepalert.AttrContext{deepalert.CtxRemote},
					},
				},
			},
			{
				Detector:    "blue",
				RuleName:    "orange",
				AlertKey:    "five",
				Description: "timeless",
				Timestamp:   time.Now(),
				Attributes: []deepalert.Attribute{
					{
						Type:    deepalert.TypeIPAddr,
						Key:     "source",
						Value:   "192.168.0.1",
						Context: []deepalert.AttrContext{deepalert.CtxRemote},
					},
				},
			},
		},
		Sections: []deepalert.ReportSection{
			{
				Author: "Familiar1",
				Attribute: deepalert.Attribute{
					Type:    deepalert.TypeIPAddr,
					Key:     "source",
					Value:   "192.168.0.1",
					Context: []deepalert.AttrContext{deepalert.CtxRemote},
				},
				Type: deepalert.ContentHost,
				Content: deepalert.ReportHost{
					RelatedDomains: []deepalert.EntityDomain{
						{
							Name:      "example.com",
							Timestamp: time.Now(),
							Source:    "tester",
						},
					},
				},
			},
			{
				Author: "Familiar2",
				Attribute: deepalert.Attribute{
					Type:    deepalert.TypeIPAddr,
					Key:     "source",
					Value:   "192.168.0.1",
					Context: []deepalert.AttrContext{deepalert.CtxRemote},
				},
				Type: deepalert.ContentHost,
				Content: deepalert.ReportHost{
					IPAddr: []string{"10.0.1.2"},
					RelatedDomains: []deepalert.EntityDomain{
						{
							Name:      "example.net",
							Timestamp: time.Now(),
							Source:    "tester",
						},
					},
				},
			},
			{
				Author: "SomeVirusScanner",
				Type:   deepalert.ContentHost,
				Attribute: deepalert.Attribute{
					Type:    deepalert.TypeIPAddr,
					Key:     "source",
					Value:   "192.168.0.2",
					Context: []deepalert.AttrContext{deepalert.CtxRemote},
				},
				Content: deepalert.ReportHost{
					RelatedMalware: []deepalert.EntityMalware{
						{
							SHA256:    "abcdefg",
							Timestamp: time.Now(),
							Scans: []deepalert.EntityMalwareScan{
								{
									Vendor: "normalVender",
									Name:   "some_malware",
								},
								{
									Vendor: "superVender",
									Name:   "some_malware2",
								},
							},
						},
					},
				},
			},
			{
				Author: "xxxx",
				Type:   deepalert.ContentUser,
				Attribute: deepalert.Attribute{
					Type:    deepalert.TypeUserName,
					Key:     "name",
					Value:   "blue",
					Context: []deepalert.AttrContext{deepalert.CtxRemote},
				},
				Content: deepalert.ReportUser{
					Activities: []deepalert.EntityActivity{
						{
							ServiceName: "magic",
							RemoteAddr:  "10.2.3.4",
						},
					},
				},
			},
		},
	}

	buf, err := main.ReportToBody(report)
	require.NotNil(t, buf)
	require.NoError(t, err)

	txt := buf.String()
	if os.Getenv("VERBOSE") != "" {
		fmt.Println(txt)
	}

	assert.Contains(t, txt, "Detected by  `blue`")
	assert.Contains(t, txt, "- source ( `ipaddr` ):  `192.168.0.1` \n")
	assert.NotContains(t, txt, "- source ( `ipaddr` ):  `192.168.0.1` \n- source ( `ipaddr` ):  `192.168.0.1`")
}
