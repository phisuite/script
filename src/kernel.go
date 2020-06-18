package main

import (
	"github.com/phisuite/schema.go"
	"strings"
)

const (
	space   = " "
	indent  = space + space
	newline = "\n"
)

type kernel struct {
	events    map[string]*schema.Event
	entities  map[string]*schema.Entity
	processes map[string]*schema.Process
}

type header struct {
	kind, name, version string
	status              schema.Status
}

type body struct {
	fields []*schema.Field
}

type behaviour struct {
	definition *schema.Process_Definition
}

type inout struct {
	kind, name, version string
}

func (k *kernel) init() {
	k.events = map[string]*schema.Event{}
	k.entities = map[string]*schema.Entity{}
	k.processes = map[string]*schema.Process{}
}

func (k *kernel) String() string {
	var str strings.Builder
	for _, event := range k.events {
		header := &header{
			kind:    "event",
			name:    event.Name,
			version: event.Version,
			status:  event.Status,
		}
		str.WriteString(header.String())
		body := &body{fields: event.Payload}
		str.WriteString(body.String())
		str.WriteString(newline)
	}
	for _, entity := range k.entities {
		header := &header{
			kind:    "entity",
			name:    entity.Name,
			version: entity.Version,
			status:  entity.Status,
		}
		str.WriteString(header.String())
		body := &body{fields: entity.Data}
		str.WriteString(body.String())
		str.WriteString(newline)
	}
	for _, process := range k.processes {
		header := &header{
			kind:    "process",
			name:    process.Name,
			version: process.Version,
			status:  process.Status,
		}
		str.WriteString(header.String())
		behaviour := &behaviour{definition: process.Definition}
		str.WriteString(behaviour.String())
		str.WriteString(newline)
	}
	return str.String()
}

func (h *header) String() string {
	var str strings.Builder
	str.WriteString(h.kind)
	switch h.status {
	case schema.Status_UNACTIVATED:
		str.WriteString("?")
	case schema.Status_ACTIVATED:
		str.WriteString("!")
	case schema.Status_DEACTIVATED:
		str.WriteString("~")
	}
	str.WriteString(space)
	str.WriteString(h.name + ":" + h.version)
	str.WriteString(newline)
	return str.String()
}

func (b *body) String() string {
	var str strings.Builder
	for _, field := range b.fields {
		str.WriteString(indent)
		str.WriteString(strings.ToLower(field.Type.String()))
		if field.Category == schema.Field_OPTIONAL {
			str.WriteString("?")
		}
		str.WriteString(space + field.Name + newline)
	}
	return str.String()
}

func (b *behaviour) String() string {
	var str strings.Builder
	inputEvent := b.definition.Input.Event
	outputEvent := b.definition.Output.Event
	errorEvent := b.definition.Error.Event
	inputEntity := b.definition.Input.Entity
	outputEntity := b.definition.Output.Entity
	errorEntity := b.definition.Error.Entity
	ios := []*inout{
		{kind: "input", name: inputEvent.Name, version: inputEvent.Version},
		{kind: "output", name: outputEvent.Name, version: outputEvent.Version},
		{kind: "error", name: errorEvent.Name, version: errorEvent.Version},
	}
	if inputEntity != nil {
		io := &inout{kind: "input", name: inputEntity.Name, version: inputEntity.Version}
		ios = append(ios, io)
	}
	if outputEntity != nil {
		io := &inout{kind: "output", name: outputEntity.Name, version: outputEntity.Version}
		ios = append(ios, io)
	}
	if errorEntity != nil {
		io := &inout{kind: "error", name: errorEntity.Name, version: errorEntity.Version}
		ios = append(ios, io)
	}
	for _, io := range ios {
		str.WriteString(io.String())
	}
	return str.String()
}

func (io *inout) String() string {
	var str strings.Builder
	str.WriteString(indent + io.kind + space)
	str.WriteString(io.name + ":" + io.version)
	str.WriteString(newline)
	return str.String()
}
