package client

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"gopkg.in/jcmturner/gokrb5.v6/iana/nametype"
	"gopkg.in/jcmturner/gokrb5.v6/krberror"
	"gopkg.in/jcmturner/gokrb5.v6/messages"
	"gopkg.in/jcmturner/gokrb5.v6/types"
)

// sessions hold TGTs and are keyed on the realm name
type sessions struct {
	Entries map[string]*session
	mux     sync.RWMutex
}

// destroy erases all sessions
func (s *sessions) destroy() {
	s.mux.Lock()
	defer s.mux.Unlock()
	for k, e := range s.Entries {
		e.destroy()
		delete(s.Entries, k)
	}
}

// update replaces a session with the one provided or adds it as a new one
func (s *sessions) update(sess *session) {
	s.mux.Lock()
	defer s.mux.Unlock()
	// if a session already exists for this, cancel its auto renew.
	if i, ok := s.Entries[sess.realm]; ok {
		if i != sess {
			// Session in the sessions cache is not the same as one provided.
			// Cancel the one in the cache and add this one.
			i.mux.Lock()
			defer i.mux.Unlock()
			i.cancel <- true
			s.Entries[sess.realm] = sess
			return
		}
	}
	// No session for this realm was found so just add it
	s.Entries[sess.realm] = sess
}

// get returns the session for the realm specified
func (s *sessions) get(realm string) (*session, bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	sess, ok := s.Entries[realm]
	return sess, ok
}

// session holds the TGT details for a realm
type session struct {
	realm                string
	authTime             time.Time
	endTime              time.Time
	renewTill            time.Time
	tgt                  messages.Ticket
	sessionKey           types.EncryptionKey
	sessionKeyExpiration time.Time
	cancel               chan bool
	mux                  sync.RWMutex
}

// AddSession adds a session for a realm with a TGT to the client's session cache.
// A goroutine is started to automatically renew the TGT before expiry.
func (cl *Client) AddSession(tgt messages.Ticket, dep messages.EncKDCRepPart) {
	if strings.ToLower(tgt.SName.NameString[0]) != "krbtgt" {
		// Not a TGT
		return
	}
	s := &session{
		realm:                tgt.SName.NameString[len(tgt.SName.NameString)-1],
		authTime:             dep.AuthTime,
		endTime:              dep.EndTime,
		renewTill:            dep.RenewTill,
		tgt:                  tgt,
		sessionKey:           dep.Key,
		sessionKeyExpiration: dep.KeyExpiration,
	}
	cl.sessions.update(s)
	cl.enableAutoSessionRenewal(s)
}

// update overwrites the session details with those from the TGT and decrypted encPart
func (s *session) update(tgt messages.Ticket, dep messages.EncKDCRepPart) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.authTime = dep.AuthTime
	s.endTime = dep.EndTime
	s.renewTill = dep.RenewTill
	s.tgt = tgt
	s.sessionKey = dep.Key
	s.sessionKeyExpiration = dep.KeyExpiration
}

// destroy will cancel any auto renewal of the session and set the expiration times to the current time
func (s *session) destroy() {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.cancel != nil {
		s.cancel <- true
	}
	s.endTime = time.Now().UTC()
	s.renewTill = s.endTime
	s.sessionKeyExpiration = s.endTime
}

// valid informs if the TGT is still within the valid time window
func (s *session) valid() bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	t := time.Now().UTC()
	if t.Before(s.endTime) && s.authTime.Before(t) {
		return true
	}
	return false
}

// tgtDetails is a thread safe way to get the session's realm, TGT and session key values
func (s *session) tgtDetails() (string, messages.Ticket, types.EncryptionKey) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.realm, s.tgt, s.sessionKey
}

// timeDetails is a thread safe way to get the session's validity time values
func (s *session) timeDetails() (string, time.Time, time.Time, time.Time, time.Time) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.realm, s.authTime, s.endTime, s.renewTill, s.sessionKeyExpiration
}

// enableAutoSessionRenewal turns on the automatic renewal for the client's TGT session.
func (cl *Client) enableAutoSessionRenewal(s *session) {
	var timer *time.Timer
	s.mux.Lock()
	s.cancel = make(chan bool, 1)
	s.mux.Unlock()
	go func(s *session) {
		for {
			s.mux.RLock()
			w := (s.endTime.Sub(time.Now().UTC()) * 5) / 6
			s.mux.RUnlock()
			if w < 0 {
				return
			}
			timer = time.NewTimer(w)
			select {
			case <-timer.C:
				renewal, err := cl.refreshSession(s)
				if !renewal && err == nil {
					// end this goroutine as there will have been a new login and new auto renewal goroutine created.
					return
				}
			case <-s.cancel:
				// cancel has been called. Stop the timer and exit.
				timer.Stop()
				return
			}
		}
	}(s)
}

// renewTGT renews the client's TGT session.
func (cl *Client) renewTGT(s *session) error {
	realm, tgt, skey := s.tgtDetails()
	spn := types.PrincipalName{
		NameType:   nametype.KRB_NT_SRV_INST,
		NameString: []string{"krbtgt", realm},
	}
	_, tgsRep, err := cl.TGSExchange(spn, cl.Credentials.Realm, tgt, skey, true, 0)
	if err != nil {
		return krberror.Errorf(err, krberror.KRBMsgError, "error renewing TGT")
	}
	s.update(tgsRep.Ticket, tgsRep.DecryptedEncPart)
	cl.sessions.update(s)
	return nil
}

// refreshSession updates either through renewal or creating a new login.
// The boolean indicates if the update was a renewal.
func (cl *Client) refreshSession(s *session) (bool, error) {
	s.mux.RLock()
	realm := s.realm
	renewTill := s.renewTill
	s.mux.RUnlock()
	if time.Now().UTC().Before(renewTill) {
		err := cl.renewTGT(s)
		return true, err
	}
	err := cl.realmLogin(realm)
	return false, err
}

// ensureValidSession makes sure there is a valid session for the realm
func (cl *Client) ensureValidSession(realm string) error {
	s, ok := cl.sessions.get(realm)
	if ok {
		s.mux.RLock()
		defer s.mux.RUnlock()
		d := s.endTime.Sub(s.authTime) / 6
		if s.endTime.Sub(time.Now().UTC()) > d {
			return nil
		}
		_, err := cl.refreshSession(s)
		return err
	}
	return cl.realmLogin(realm)
}

// sessionTGTDetails is a thread safe way to get the TGT and session key values for a realm
func (cl *Client) sessionTGT(realm string) (tgt messages.Ticket, sessionKey types.EncryptionKey, err error) {
	err = cl.ensureValidSession(realm)
	if err != nil {
		return
	}
	s, ok := cl.sessions.get(realm)
	if !ok {
		err = fmt.Errorf("could not find TGT session for %s", realm)
		return
	}
	_, tgt, sessionKey = s.tgtDetails()
	return
}

func (cl *Client) sessionTimes(realm string) (authTime, endTime, renewTime, sessionExp time.Time, err error) {
	s, ok := cl.sessions.get(realm)
	if !ok {
		err = fmt.Errorf("could not find TGT session for %s", realm)
		return
	}
	_, authTime, endTime, renewTime, sessionExp = s.timeDetails()
	return
}

// spnRealm resolves the realm name of a service principal name
func (cl *Client) spnRealm(spn types.PrincipalName) string {
	return cl.Config.ResolveRealm(spn.NameString[len(spn.NameString)-1])
}