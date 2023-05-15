// ca is a convenience package providing CodeComet root CA as an embed
// The root is NOT trusted by default.
// To trust it:
// - myConfig.Server.ClientCA = ca.CodeComet // uses CodeComet CA for client cert validation
// - myConfig.Client.RootCAs = []string{ca.CodeComet} // uses CodeComet CA as an extra root CA for cert validation
// This is also done if passing true to go-core/config.New(true, "myconfiglocation")
package ca

import (
	_ "embed"
)

//go:embed codecomet.crt
var CodeComet string
