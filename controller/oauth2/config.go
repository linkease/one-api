package oauth2

import "net/url"

var (
	client_id     = "aitest"
	client_secret = "aitest123456"
	oauth2_url    = url.URL{
		Scheme: "https",
		Host:   "sso.koolcenter.com",
	}
	callback_url = url.URL{
		Scheme: "http",
		Host:   "localhost:3000",
		Path:   "/auth/oauth2_basic/callback",
	}
)

// http://localhost:3000/auth/oauth2_basic
// http://localhost:9096/oauth/authorize

// clientId, clientSecret oauth2的域名 本机域名
func SetConfig(clientId, clientSecret, oauth2Url, localDomain string) error {
	client_id = clientId
	client_secret = clientSecret
	uri1, err := url.Parse(oauth2Url)
	if err != nil {
		return err
	}
	oauth2_url.Scheme = uri1.Scheme
	oauth2_url.Host = uri1.Host

	uri2, err := url.Parse(localDomain)
	if err != nil {
		return err
	}
	callback_url.Scheme = uri2.Scheme
	callback_url.Host = uri2.Host
	return nil
}
