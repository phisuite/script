package main

import (
	"context"
	"github.com/phisuite/schema.go"
	"google.golang.org/grpc"
)

type updater struct {
	kernel *kernel
	parser *codeParser
	inspector *inspector
}

func (u *updater) update(filename, editorAddr, inspectorAddr string) (err error) {
	if inspectorAddr == "" {
		inspectorAddr = editorAddr
	}
	u.kernel = &kernel{}
	u.kernel.init()
	u.parser = &codeParser{}
	u.parser.Parse(filename)
	u.inspector = &inspector{}
	err = u.inspector.inspect(inspectorAddr)
	if err != nil {
		return
	}

	conn, err := grpc.Dial(editorAddr, grpc.WithInsecure())
	if err != nil {
		return
	}
	defer conn.Close()

	err = u.updateEvents(conn)
	if err != nil {
		return
	}
	err = u.updateEntities(conn)
	if err != nil {
		return
	}
	err = u.updateProcesses(conn)
	return
}

func (u *updater) updateEvents(conn *grpc.ClientConn) (err error) {
	writer := schema.NewEventWriteAPIClient(conn)
	for id, event := range u.parser.kernel.events {
		options := &schema.Options{Name: event.Name, Version: event.Version}
		existing, ok := u.inspector.kernel.events[id]
		if !ok {
			existing, err = writer.Create(context.Background(), event)
		} else if existing.Status == schema.Status_UNACTIVATED {
			existing, err = writer.Update(context.Background(), event)
		}
		if err != nil {
			return
		}
		var result *schema.Event
		switch event.Status {
		case existing.Status:
		case schema.Status_ACTIVATED:
			result, err = writer.Activate(context.Background(), options)
		case schema.Status_DEACTIVATED:
			result, err = writer.Deactivate(context.Background(), options)
		default:
		}
		if err != nil {
			return
		}
		if result != nil {
			existing.Status = result.Status
		}
		u.kernel.events[id] = existing
	}
	return
}

func (u *updater) updateEntities(conn *grpc.ClientConn) (err error) {
	writer := schema.NewEntityWriteAPIClient(conn)
	for id, entity := range u.parser.kernel.entities {
		options := &schema.Options{Name: entity.Name, Version: entity.Version}
		existing, ok := u.inspector.kernel.entities[id]
		if !ok {
			existing, err = writer.Create(context.Background(), entity)
		} else if existing.Status == schema.Status_UNACTIVATED {
			existing, err = writer.Update(context.Background(), entity)
		}
		if err != nil {
			return
		}
		var result *schema.Entity
		switch entity.Status {
		case existing.Status:
		case schema.Status_ACTIVATED:
			result, err = writer.Activate(context.Background(), options)
		case schema.Status_DEACTIVATED:
			result, err = writer.Deactivate(context.Background(), options)
		default:
		}
		if err != nil {
			return
		}
		if result != nil {
			existing.Status = result.Status
		}
		u.kernel.entities[id] = existing
	}
	return
}

func (u *updater) updateProcesses(conn *grpc.ClientConn) (err error) {
	writer := schema.NewProcessWriteAPIClient(conn)
	for id, process := range u.parser.kernel.processes {
		options := &schema.Options{Name: process.Name, Version: process.Version}
		existing, ok := u.inspector.kernel.processes[id]
		if !ok {
			existing, err = writer.Create(context.Background(), process)
		} else if existing.Status == schema.Status_UNACTIVATED {
			existing, err = writer.Update(context.Background(), process)
		}
		if err != nil {
			return
		}
		var result *schema.Process
		switch process.Status {
		case existing.Status:
		case schema.Status_ACTIVATED:
			result, err = writer.Activate(context.Background(), options)
		case schema.Status_DEACTIVATED:
			result, err = writer.Deactivate(context.Background(), options)
		default:
		}
		if err != nil {
			return
		}
		if result != nil {
			existing.Status = result.Status
		}
		u.kernel.processes[id] = existing
	}
	return
}
