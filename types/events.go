// Copyright (c) 2024 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package types contains types by the whatsmeow libraryn// MessageSource the source information of a group, this is the group membert// Whether the message was sent by the current user instead of another contact.
	FromMe bool
	// Whether the message was sent in a group chat.
	IsGroup bool
}

// MessageInfo contains metadata about an incoming message.
type MessageInfo struct {
	MessageSource
	// The unique message ID.
	ID MessageID
	// The server-assigned timestamp of the message.
	Timestamp time.Time
	// The push name (display name) of the sender.
	PushName string
	// Whether the message was broadcast (e.g. status update).
	Broadcast bool
}

// MessageID is the ID of a WhatsApp message.
type MessageID = string

// ReceiptType represents the type of a message receipt.
type ReceiptType string

const (
	// ReceiptTypeDelivered indicates the message was delivered to the device.
	ReceiptTypeDelivered ReceiptType = ""
	// ReceiptTypeRead indicates the message was read by the user.
	ReceiptTypeRead ReceiptType = "read"
	// ReceiptTypePlayed indicates the media message was played by the user.
	ReceiptTypePlayed ReceiptType = "played"
)

// Receipt is emitted when a receipt is received for a sent message.
type Receipt struct {
	// The source of the receipt (who sent it and in which chat).
	MessageSource
	// The IDs of the messages this receipt is for.
	MessageIDs []MessageID
	// The timestamp of the receipt.
	Timestamp time.Time
	// The type of the receipt.
	Type ReceiptType
}

// Presence represents a contact's presence state.
type Presence string

const (
	// PresenceAvailable means the contact is online.
	PresenceAvailable Presence = "available"
	// PresenceUnavailable means the contact is offline.
	PresenceUnavailable Presence = "unavailable"
)

// PresenceEvent is emitted when a contact's presence changes.
type PresenceEvent struct {
	// The JID of the contact.
	From JID
	// The new presence state. True means the contact is offline/unavailable.
	Unavailable bool
	// The last seen time (only set when Unavailable is true).
	// Note: WhatsApp only provides this if the contact has their last seen visible to you.
	LastSeen time.Time
}

// Connected is emitted when the client successfully connects to WhatsApp.
type Connected struct{}

// Disconnected is emitted when the client disconnects from WhatsApp.
type Disconnected struct {
	// Whether the disconnect was requested by the user (i.e. logged out).
	// If false, the disconnect was likely due to a network error or server-side kick.
	LoggedOut bool
}

// QR is emitted when a QR code is available for scanning.
type QR struct {
	// The QR code codes to display. The client should cycle
