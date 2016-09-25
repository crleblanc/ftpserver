package server

import "fmt"

func (p *Paradise) HandleUser() {
	p.user = p.param
	p.writeMessage(331, "User name ok, password required")
}

func (p *Paradise) HandlePass() {
	// think about using https://developer.bitium.com
	if AuthManager.CheckUser(p.user, p.param, &p.userInfo) {
		p.writeMessage(230, "Password ok, continue")
	} else {
		p.writeMessage(530, "Incorrect password, not logged in")
		p.theConnection.Close()
		delete(ConnectionMap, p.cid)
	}
}

func (p *Paradise) HandleAuth() {
	fmt.Println(p.param)
	p.writeMessage(234, "AUTH command ok. Expecting TLS Negotiation.")
	//2012-02-08 08:04:11 7424 3 Status: Initializing TLS...
	//2012-02-08 08:04:11 7424 3 Status: Verifying certificate...
	//2012-02-08 08:04:11 7424 3 Command: USER msoftpuser
	//2012-02-08 08:04:11 7424 3 Status: TLS/SSL connection established.
}

func (p *Paradise) HandlePbsz() {
	fmt.Println(p.param)
	p.writeMessage(200, "PBSZ OK.")
}

func (p *Paradise) HandleProt() {
	fmt.Println(p.param)
	p.writeMessage(200, "PROT OK.")
}
