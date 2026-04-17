// Copyright (c) 2024 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package whatsbn	"sync"
t/.mautypes/events// EventHandler is a function that can WhatsApp client struct.
type Client struct {
	Store   *store.Device
	Log     log.Logger
	
	// Event handlers
	eEventHandler struct {
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
// Each handler call is wrapped in a recover() so a panicking handler does not
// prevent the remaining handlers from receiving the event.
//
// NOTE(personal): switched from sequential-fail-on-panic to recover-per-handler.
// This makes the behavior more robust for my use case where multiple independent
// handlers are registered and a bug in one shouldn't silently drop events for others.
//
// NOTE(personal): handlers are copied under the read lock so we don't hold the
// lock while invoking user code, which could deadlock if a handler calls
// AddEventHandler or RemoveEventHandler.
//
// NOTE(personal): added log output on panic so panics are visible in logs rather
// than silently swallowed.
func (cli *Client) dispatchEvent(evt interface{}) {
	cli.eventHandlersLock.RLock()
	handlers := cli.eventHandlers
	cli.eventHandlersLock.RUnlock()
	for _, handler := range handlers {
		fun
