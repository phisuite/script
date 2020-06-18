package main

import (
	"bufio"
	"github.com/phisuite/schema.go"
	"log"
	"os"
	"strings"
)

type codeParser struct {
	kernel *kernel
	parser FileParser
}

func (c *codeParser) Parse(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	c.parser = &fileParser{
		scanner: bufio.NewScanner(file),
	}
	c.kernel = &kernel{}
	c.kernel.init()
	c.parser.Next()
	for c.parser.HasNext() {
		token := c.parser.Token()
		switch {
		case token.Is("event"):
			c.extractEvent()
		case token.Is("entity"):
			c.extractEntity()
		case token.Is("process"):
			c.extractProcess()
		default:
			c.parser.Fatal("invalid token %s", token)
		}
	}
}

func (c *codeParser) extractEvent() {
	event := &schema.Event{}
	event.Status = c.extractStatus()
	c.parser.Next()
	event.Name, event.Version = c.extractId()
	c.kernel.events[c.parser.Token().String()] = event
	c.parser.Next()
	event.Payload = c.extractFields()
}

func (c *codeParser) extractEntity() {
	entity := &schema.Entity{}
	entity.Status = c.extractStatus()
	c.parser.Next()
	entity.Name, entity.Version = c.extractId()
	c.kernel.entities[c.parser.Token().String()] = entity
	c.parser.Next()
	entity.Data = c.extractFields()
}

func (c *codeParser) extractStatus() schema.Status {
	statusToken := c.parser.Token().Last()
	switch statusToken {
	case '?':
		return schema.Status_UNACTIVATED
	case '!':
		return schema.Status_ACTIVATED
	case '~':
		return schema.Status_DEACTIVATED
	default:
		c.parser.Fatal("invalid status %c", statusToken)
	}
	return 0
}

func (c *codeParser) extractId() (string, string) {
	id := c.parser.Token().Split(":")
	if len(id) != 2 {
		c.parser.Fatal("invalid id %s", c.parser.Token())
	}
	return id[0], id[1]
}

func (c *codeParser) extractFields() []*schema.Field {
	if c.parser.Token().Is("event", "entity", "process") {
		return []*schema.Field{}
	}
	field := &schema.Field{}
	fieldProps := c.parser.Token().Split("?")
	fieldTypeValue, ok := schema.Field_Type_value[strings.ToUpper(fieldProps[0])]
	if !ok {
		c.parser.Fatal("invalid type %s", fieldProps[0])
	}
	field.Type = schema.Field_Type(fieldTypeValue)
	field.Category = schema.Field_REQUIRED
	if len(fieldProps) == 2 {
		field.Category = schema.Field_OPTIONAL
	}
	c.parser.Next()
	field.Name = c.parser.Token().String()
	c.parser.Next()
	return append(c.extractFields(), field)
}

func (c *codeParser) extractProcess() {
	process := &schema.Process{}
	process.Status = c.extractStatus()
	c.parser.Next()
	process.Name, process.Version = c.extractId()
	c.kernel.processes[c.parser.Token().String()] = process
	c.parser.Next()
	definition := &schema.Process_Definition{
		Input:  &schema.Process_Data{},
		Output: &schema.Process_Data{},
		Error:  &schema.Process_Data{},
	}
	for c.parser.HasNext() {
		switch {
		case c.parser.Token().Is("input"):
			c.extractProcessData(definition.Input)
		case c.parser.Token().Is("output"):
			c.extractProcessData(definition.Output)
		case c.parser.Token().Is("error"):
			c.extractProcessData(definition.Error)
		default:
		}
	}
	process.Definition = definition
}

func (c *codeParser) extractProcessData(processData *schema.Process_Data) {
	c.parser.Next()
	name, version := c.extractId()
	id := name + ":" + version
	if event, ok := c.kernel.events[id]; ok {
		processData.Event = event
		c.parser.Next()
		return
	}
	if entity, ok := c.kernel.entities[id]; ok {
		processData.Entity = entity
		c.parser.Next()
		return
	}
	c.parser.Fatal("%s not defined", id)
}
