// ca is a convenience package providing the ability to embed your root CA
// This root embed is **NOT** trusted by default.
// To trust it:
// - myConfig.Server.ClientCA = ca.MyRoot // uses CA for client cert validation
// - myConfig.Client.RootCAs = []string{ca.MyRoot} // uses CA as an extra root CA for cert validation
// This above is also done if passing true to go-core/config.New(true, "myconfiglocation")
// Right now, the cert is invalid and zeroed-out
package ca

import (
	_ "embed"
)

//go:embed myroot.crt
var MyRoot string
