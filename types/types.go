// Copyright (c) 2024 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package types contains types used by WhatsApp JID ( JID struct {
	  uint8
	Device uint8
	Server string
	AD     bool
}

// Known WhatsApp servers.
const (
	DefaultUserServer = "s.whatsapp.net"
	GroupServer       = "g.us"
	LegacyUserServer  = "c.us"
	BroadcastServer   = "broadcast"
	MessengerServer   = "msgr"
	InteropServer     = "interop"
)

// NewJID creates a new JID with the given user and server.
func NewJID(user, server string) JID {
	return JID{User: user, Server: server}
}

// IsEmpty returns true if the JID has no user and no server.
func (jid JID) IsEmpty() bool {
	return jid.User == "" && jid.Server == ""
}

// String returns the string representation of the JID.
func (jid JID) String() string {
	if jid.AD {
		return fmt.Sprintf("%s.%d:%d@%s", jid.User, jid.Agent, jid.Device, jid.Server)
	}
	if jid.User == "" {
		return jid.Server
	}
	return fmt.Sprintf("%s@%s", jid.User, jid.Server)
}

// ParseJID parses a JID from a string.
func ParseJID(s string) (JID, error) {
	if s == "" {
		return JID{}, fmt.Errorf("empty JID")
	}
	atIdx := strings.LastIndex(s, "@")
	if atIdx < 0 {
		return NewJID("", s), nil
	}
	server := s[atIdx+1:]
	user := s[:atIdx]
	if colonIdx := strings.Index(user, ":"); colonIdx > 0 && strings.Contains(user, ".") {
		var agent, device uint8
		dotIdx := strings.Index(user, ".")
		_, err := fmt.Sscanf(user[dotIdx+1:], "%d:%d", &agent, &device)
		if err == nil {
			return JID{
				User:   user[:dotIdx],
				Agent:  agent,
				Device: device,
				Server: server,
				AD:     true,
			}, nil
		}
	}
	return NewJID(user, server), nil
}

// MessageID is the ID of a WhatsApp message.
type MessageID = string

// MessageSource contains the basic sender and chat information of a message.
type MessageSource struct {
	Chat     JID
	Sender   JID
	IsFromMe bool
	IsGroup  bool
}

// UserInfo contains info about a WhatsApp user.
type UserInfo struct {
	VerifiedName string
	Status       string
	PictureID    string
	Devices      []JID
}
