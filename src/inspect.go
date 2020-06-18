package main

import (
	"context"
	"github.com/phisuite/schema.go"
	"google.golang.org/grpc"
	"io"
)

type inspector struct {
	kernel *kernel
}

func (i *inspector) inspect(inspectorAddr string) (err error) {
	i.kernel = &kernel{}
	i.kernel.init()

	conn, err := grpc.Dial(inspectorAddr, grpc.WithInsecure())
	if err != nil {
		return
	}
	defer conn.Close()

	err = i.inspectEvents(conn)
	if err != nil {
		return
	}
	err = i.inspectEntities(conn)
	if err != nil {
		return
	}
	err = i.inspectProcesses(conn)
	return
}

func (i *inspector) inspectEvents(conn *grpc.ClientConn) (err error) {
	reader := schema.NewEventReadAPIClient(conn)

	stream, err := reader.List(context.Background(), &schema.Options{})
	if err != nil {
		return
	}
	for {
		event, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if event == nil {
			continue
		}
		id := event.Name+":"+event.Version
		i.kernel.events[id] = event
	}
	return nil
}

func (i *inspector) inspectEntities(conn *grpc.ClientConn) error {
	reader := schema.NewEntityReadAPIClient(conn)

	stream, err := reader.List(context.Background(), &schema.Options{})
	if err != nil {
		return err
	}
	for {
		entity, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if entity == nil {
			continue
		}
		id := entity.Name+":"+ entity.Version
		i.kernel.entities[id] = entity
	}
	return nil
}

func (i *inspector) inspectProcesses(conn *grpc.ClientConn) error {
	reader := schema.NewProcessReadAPIClient(conn)

	stream, err := reader.List(context.Background(), &schema.Options{})
	if err != nil {
		return err
	}
	for {
		process, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if process == nil {
			continue
		}
		id := process.Name+":"+ process.Version
		i.kernel.processes[id] = process
	}
	return nil
}
