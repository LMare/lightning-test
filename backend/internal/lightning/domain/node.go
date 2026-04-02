package domain

type NodeConfigDescriptor struct {
	AuthData	LndClientAuthData	`yaml:"data"`
	Id			string					`yaml:"id"`
}


type LndClientAuthData struct {
	TlsCertPath 	string	`yaml:"cert"`
	MacaroonPath 	string	`yaml:"macaroon"`
	LndAddress 		string	`yaml:"url"`
}

func NewLndClientAuthData(c, m, a string) LndClientAuthData {
	return LndClientAuthData{c, m, a}
}
