package connections

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

type Connection struct {
	uuid       uuid.UUID
	response   http.ResponseWriter
	request    *http.Request
	ws         websocket.Upgrader
	connection *websocket.Conn
	stop       chan bool
	reload     chan bool
	fail       chan error
}

func (c *Connection) UUID() string {
	return c.uuid.String()
}

func (c *Connection) Open() error {
	c.ws = websocket.Upgrader{}
	c.ws.CheckOrigin = func(r *http.Request) bool {
		origin := r.Header.Values("Origin")
		if len(origin) != 1 {
			return false
		}

		return strings.HasPrefix(origin[0], "http://localhost")
	}

	connection, err := c.ws.Upgrade(c.response, c.request, nil)
	if err != nil {
		return fmt.Errorf("(%s) failed to upgrade the connection", c.UUID())
	}
	c.connection = connection

	return nil
}

func (c *Connection) Start() error {

	c.stop = make(chan bool)
	c.reload = make(chan bool)
	c.fail = make(chan error)

	go c.alive()
	go c.watch()

	select {
	case <-c.stop:
		{
			c.stop = nil
			c.reload = nil
			c.fail = nil

			return nil
		}
	case err := <-c.fail:
		{
			return fmt.Errorf("(%s) something was bad while the connection was active: %v", c.UUID(), err)
		}
	}
}

func (c *Connection) Reload() error {
	if c.reload == nil {
		return fmt.Errorf("(%s) failed to reload, so the connection is not started", c.UUID())
	}
	c.reload <- true
	return nil
}

func (c *Connection) Stop() error {
	if c.stop == nil {
		return fmt.Errorf("(%s) failed to stop, so the connection is not started", c.UUID())
	}

	c.stop <- true
	return nil
}

func (c *Connection) Close() error {
	err := c.connection.Close()
	if err != nil {
		return fmt.Errorf("(%s) failed to close connection: %v", c.UUID(), err)
	}
	return nil
}

func (c *Connection) alive() {
	for {
		_, _, err := c.connection.ReadMessage()
		if err != nil {
			c.stop <- true
		}
	}
}

func (c *Connection) watch() {
	for {
		select {
		case <-c.reload:
			{
				err := c.connection.WriteMessage(1, []byte("reload"))
				if err != nil {
					log.Printf("(%s) failed to send reload signal: %s", c.UUID(), err)
					break
				}
			}
		case <-c.stop:
			{
				log.Printf("(%s) stoping watcher", c.UUID())
				return
			}
		}
	}
}

type Configurer interface {
	ResponseWriter(response http.ResponseWriter) error
	Request(request *http.Request) error
}

type configurerPool struct {
	pool []func(c *Connection) error
}

func (cp *configurerPool) ResponseWriter(response http.ResponseWriter) error {

	cp.pool = append(cp.pool, func(c *Connection) error {
		c.response = response
		return nil
	})

	return nil
}

func (cp *configurerPool) Request(request *http.Request) error {

	cp.pool = append(cp.pool, func(c *Connection) error {
		c.request = request
		return nil
	})

	return nil
}

func New(options ...func(Configurer) error) (*Connection, error) {

	configuration := &configurerPool{}
	conn := &Connection{
		uuid: uuid.New(),
	}

	for _, option := range options {
		err := option(configuration)
		if err != nil {
			return nil, fmt.Errorf("(%s) failed to prepare the configuration: %v", conn.UUID(), err)
		}
	}

	for _, config := range configuration.pool {
		err := config(conn)
		if err != nil {
			return nil, fmt.Errorf("(%s) failed to apply the configuration: %v", conn.UUID(), err)
		}
	}

	if conn.request == nil {
		return nil, fmt.Errorf("(%s) request is required", conn.UUID())
	}

	if conn.response == nil {
		return nil, fmt.Errorf("(%s) response is required", conn.UUID())
	}

	return conn, nil
}