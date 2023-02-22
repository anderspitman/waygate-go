package waygate

func NewCaddyLayer4Config() *CaddyLayer4Config {
	c := &CaddyLayer4Config{
		Apps: CaddyLayer4Apps{
			Layer4: CaddyLayer4{
				Servers: CaddyServers{
					Waygate: CaddyWaygate{
						Listen: []string{
							":443",
						},
						Routes: []CaddyRoute{},
					},
				},
			},
		},
	}
	return c
}

func NewCaddyHttpConfig() *CaddyHttpConfig {
	c := &CaddyHttpConfig{
		Apps: CaddyHttpApps{
			Http: CaddyHttp{
				Servers: CaddyServers{
					Waygate: CaddyWaygate{
						Routes: []CaddyRoute{},
					},
				},
			},
			Tls: CaddyTls{
				Automation: CaddyAutomation{
					Policies: []CaddyPolicy{
						CaddyPolicy{
							OnDemand: true,
							Issuers: []CaddyIssuer{
								CaddyIssuer{
									Module: "acme",
									Challenges: CaddyChallenges{
										Http: CaddyHttpChallenge{
											Disabled: true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return c
}

type CaddyLayer4Config struct {
	Apps CaddyLayer4Apps `json:"apps"`
}

type CaddyLayer4Apps struct {
	Layer4 CaddyLayer4 `json:"layer4"`
}

type CaddyLayer4 struct {
	Servers CaddyServers `json:"servers"`
}

type CaddyServers struct {
	Waygate CaddyWaygate `json:"waygate"`
}

type CaddyWaygate struct {
	Listen []string     `json:"listen"`
	Routes []CaddyRoute `json:"routes"`
}

type CaddyRoute struct {
	Match  []interface{}  `json:"match"`
	Handle []CaddyHandler `json:"handle"`
}

type CaddyMatcher struct {
	Tls CaddyTlsMatcher `json:"tls"`
}

type CaddyTlsMatcher struct {
	Sni []string `json:"sni"`
}

type CaddyHostMatcher struct {
	Host []string `json:"host"`
}

type CaddyHandler struct {
	Handler   string          `json:"handler"`
	Upstreams []CaddyUpstream `json:"upstreams"`
}

type CaddyUpstream struct {
	Dial string `json:"dial"`
}

type CaddyHttpConfig struct {
	Apps CaddyHttpApps `json:"apps"`
}

type CaddyHttpApps struct {
	Http CaddyHttp `json:"http"`
	Tls  CaddyTls  `json:"tls"`
}

type CaddyHttp struct {
	Servers CaddyServers `json:"servers"`
}

type CaddyTls struct {
	Automation CaddyAutomation `json:"automation"`
}

type CaddyAutomation struct {
	Policies []CaddyPolicy `json:"policies"`
}

type CaddyPolicy struct {
	OnDemand bool          `json:"on_demand"`
	Issuers  []CaddyIssuer `json:"issuers"`
}

type CaddyIssuer struct {
	Module     string          `json:"module"`
	Challenges CaddyChallenges `json:"challenges"`
}

type CaddyChallenges struct {
	Http CaddyHttpChallenge `json:"http"`
}

type CaddyHttpChallenge struct {
	Disabled bool `json:"disabled"`
}
