// Copyright (c) 2024 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package whatsmeow implementsatsApp web client.
package whatsmeow

import (
	/atomic"

au.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"go.mau.fi/whatsmeow/util/log"
)

// EventHandler is a function that can be registered to receive events.
type EventHandler func(evt interface{})

// Client is the main WhatsApp client struct.
type Client struct {
	Store   *store.Device
	Log     log.Logger
	RecipientDevicesCache RecipientDevicesCache

	// Event handlers
	eventHandlersLock sync.RWMutex
	eventHandlers     []wrappedEventHandler
	lastHandlerID     uint32

	// Connection state
	connected     atomic.Bool
	connectCancel context.CancelFunc
	connectLock   sync.Mutex

	// Unique device identifier
	uniqueID  string
	identityID []byte
}

type wrappedEventHandler struct {
	fn EventHandler
	id uint32
}

// RecipientDevicesCache is an interface for caching recipient device lists.
type RecipientDevicesCache interface {
	GetDevices(ctx context.Context, jids []types.JID) ([]types.JID, error)
}

// NewClient creates a new WhatsApp client with the given device store and logger.
func NewClient(deviceStore *store.Device, log log.Logger) *Client {
	if log == nil {
		log = logger.Sub("Client")
	}
	return &Client{
		Store: deviceStore,
		Log:   log,
	}
}

// AddEventHandler registers a new event handler function and returns its ID.
// The ID can be used to remove the handler later with RemoveEventHandler.
func (cli *Client) AddEventHandler(handler EventHandler) uint32 {
	newID := atomic.AddUint32(&cli.lastHandlerID, 1)
	cli.eventHandlersLock.Lock()
	cli.eventHandlers = append(cli.eventHandlers, wrappedEventHandler{fn: handler, id: newID})
	cli.eventHandlersLock.Unlock()
	return newID
}

// RemoveEventHandler removes the event handler with the given ID.
// Returns true if the handler was found and removed.
func (cli *Client) RemoveEventHandler(id uint32) bool {
	cli.eventHandlersLock.Lock()
	defer cli.eventHandlersLock.Unlock()
	for i, handler := range cli.eventHandlers {
		if handler.id == id {
			cli.eventHandlers = append(cli.eventHandlers[:i], cli.eventHandlers[i+1:]...)
			return true
		}
	}
	return false
}

// RemoveAllEventHandlers removes all registered event handlers.
func (cli *Client) RemoveAllEventHandlers() {
	cli.eventHandlersLock.Lock()
	cli.eventHandlers = nil
	cli.eventHandlersLock.Unlock()
}

// dispatchEvent sends the given event to all registered event handlers.
func (cli *Client) dispatchEvent(evt interface{}) {
	cli.eventHandlersLock.RLock()
	handlers := cli.eventHandlers
	cli.eventHandlersLock.RUnlock()
	for _, handler := range handlers {
		handler.fn(evt)
	}
}

// IsConnected returns true if the client is currently connected to WhatsApp.
func (cli *Client) IsConnected() bool {
	return cli.connected.Load()
}

// IsLoggedIn returns true if the client has valid credentials stored.
func (cli *Client) IsLoggedIn() bool {
	return cli.Store != nil && cli.Store.ID != nil
}

// Disconnect disconnects the client from WhatsApp.
func (cli *Client) Disconnect() {
	cli.connectLock.Lock()
	defer cli.connectLock.Unlock()
	if cli.connectCancel != nil {
		cli.connectCancel()
		cli.connectCancel = nil
	}
	cli.connected.Store(false)
	cli.dispatchEvent(&events.Disconnected{})
}
