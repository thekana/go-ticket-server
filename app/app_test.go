package app

import (
	"github.com/spf13/viper"

	"ticket-reservation/log"
)

var (
	privateKeyBytes = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA64fKXhCFrzhFmXJO0aJvlrB/lpRZsAXjzfvLwxnC5nFEbD7H
HTom6v7Fxk+OD3yONvcGAxxcZ0gC4COiGmhQJCdmRfR/7HHcODV9jIDrWcd8CdYU
Y1JD12GzQ79OyG2QkDELSODPY0s34PGd86EyYJQHRJzdbOKcBjNcEDTUwvTy5V/w
oPncT/LLVOFfDNSPsRjb2IdK9BjZ8F0xa7nkErE/GWotH3ioQJbIDyuCw4aefqQc
iPSmRiJ8GQSWp6Js3RQZkRW5dW3OZARfS8RQUV3FKiqcrryrYxsIoTHWxcEm2bgw
emEQgJFPcVrmvzaSOyDopLUyiLA0KoWgb0DO/wIDAQABAoIBAQDGKzFcpaAlRj5h
Fik0/uvOqOAg7N7tWHdMZ3AQzosK4wBD5yoW8Eewbv4uld8cLptlqb/YPDOO/qrU
tYZ9m4jacn/9mfNMGJzGelMRNaNPnaVCjWoICz5jaKOw+7SotG7usyUgg25ax/3S
+NgCFX3SfpoC48z2AVBvcyUhyaG+CUdhk5bvMxk92MazMTLOdNYmIjA0eO4gICfj
ixeUkkJhlKF+pQIjajEqxyBp2xJW8WUAgjpXMsNZxv50HljoYk5uhtbndF/ZJ8w3
SL8ZnURqEDploQVY/pfgx24pStDODdeyAOmPsCW5MLg14AIRa3nBCMSZ0a/zfbrH
fwDXR4b5AoGBAPc+WSNOXTIeLph/5WNbPhgMuHinLSZpu3z2fvVE3Km+G4Vse4Hn
5SxwisgsLNPlziD+4cR89mYHhR9dltWYoVb22qZanv88XuZHZFmZN6kupg3sPon/
JyeZm252pfaB+A/VVCjtsuYHGMbGIjFlXw/ZZUkG6bGBGDkhXt+ODlBdAoGBAPPf
PpXuq20nNq3Q1GGl0zC4oXwr+a4Ep79fphwGuoF4aq7Y++fvrujgt50t9kY66eby
p3t33oCDzDqIrlCKbryBiJPHjxvz5M4zan7U5g0t7AgW7Ak0nVfHGyZNEIZy5WNT
CMfFUAn/ZSYin+j04pE40Jr1A9My8FohRVsuZxcLAoGAOi5Dz5kbOTX9BQnjsvO/
su0bY9kDOOzcn9VpntHrk04XL9iNX85wEXsSTXSHv/1t+jnAavp1CSFv5Cej3POj
09EXNtpQXAOa2VmndaYmgPtnPBOBy/ts/VaaSu5Es7N16lPrEA6PcK3u2Ke7WCBg
tFwWB49G4uxcBOWja7wEBkkCgYA+0BFerMqaoq8Ctfb976gltGhjgzAcEjbio9A1
B8ah8lIIFvtLEgELGlYwtdXo4OO+CGH6+zTkBQ5lRS8gr4c2Jmb3KT9DA966/aNA
Z7WZT2qr6ruA43xjT3U+uDq5Zn6OxqRMUBX9fTqgR+rIJcr1fJy+TL1feI9Pp6Il
ih4jYQKBgHn+Q+K312VL3Hhi9reLSYyLGOugPx5y2Uh2Qx0sfN+dkuw92MgBnEDT
st2wfLWAxeThIfCjozIsaXY+Yy/u2QCYljqhMcHyNDSJgBQd80qWENCE3o2D7Wk9
O4Oha6Sn/kXkezqZGWOKiLCC7oyAO7PEM+2Rp/Js+ifER6j+8Y/h
-----END RSA PRIVATE KEY-----`)
	publicKeyBytes = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA64fKXhCFrzhFmXJO0aJv
lrB/lpRZsAXjzfvLwxnC5nFEbD7HHTom6v7Fxk+OD3yONvcGAxxcZ0gC4COiGmhQ
JCdmRfR/7HHcODV9jIDrWcd8CdYUY1JD12GzQ79OyG2QkDELSODPY0s34PGd86Ey
YJQHRJzdbOKcBjNcEDTUwvTy5V/woPncT/LLVOFfDNSPsRjb2IdK9BjZ8F0xa7nk
ErE/GWotH3ioQJbIDyuCw4aefqQciPSmRiJ8GQSWp6Js3RQZkRW5dW3OZARfS8RQ
UV3FKiqcrryrYxsIoTHWxcEm2bgwemEQgJFPcVrmvzaSOyDopLUyiLA0KoWgb0DO
/wIDAQAB
-----END PUBLIC KEY-----`)
)

func getLoggerForTesting() (log.Logger, error) {
	logLevel := viper.GetString("Log.Level")
	logLevel = log.NormalizeLogLevel(logLevel)

	logColor := viper.GetBool("Log.Color")
	logJSON := viper.GetBool("Log.JSON")

	logger, err := log.NewLogger(&log.Configuration{
		EnableConsole:     true,
		ConsoleLevel:      logLevel,
		ConsoleJSONFormat: logJSON,
		Color:             logColor,
	}, log.InstanceZapLogger)
	if err != nil {
		return nil, err
	}
	return logger, nil
}
