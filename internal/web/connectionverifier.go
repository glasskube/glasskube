package web

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type connectionVerifier struct {
	enabled    bool
	mutex      sync.Mutex
	operations map[string]operation
	failed     chan struct{}
}

type operation struct {
	timestamp       time.Time
	resourceVersion string
	eventReceived   bool
}

func verifier() *connectionVerifier {
	return &connectionVerifier{
		failed:     make(chan struct{}),
		operations: make(map[string]operation),
	}
}

func (c *connectionVerifier) start() {
	c.enabled = true
	go func() {
		// TODO write logs to a logfile?

		/*
			problems with this approach:
			- it only concerns (cluster-)packages, not repos or package infos
			- it would only detect when the user actually does something, but not if connection drops randomly
		*/

		tick := time.NewTicker(2 * time.Second)
		defer tick.Stop()
		for {
			<-tick.C
			c.mutex.Lock()
			for key, op := range c.operations {
				fmt.Fprintf(os.Stderr, "VERIFIER: Checking operations for %s\n", key)
				if time.Now().After(op.timestamp.Add(5 * time.Second)) {
					if op.eventReceived {
						fmt.Fprintf(os.Stderr, "VERIFIER: update of %s probably ignorable – deleting\n", key)
						delete(c.operations, key)
					} else {
						fmt.Fprintf(os.Stderr, "VERIFIER: %s seems to be outdated – are we stuck???\n", key)
						// TODO this would be the point where we show the error to the user or try to restart everything
						c.failed <- struct{}{}
					}
				}
			}
			c.mutex.Unlock()
		}
	}()
}

func (c *connectionVerifier) expectAdd(key string) {
	if !c.enabled {
		return
	}
	// TODO
}

func (c *connectionVerifier) expectUpdate(key string, resourceVersion string) {
	if !c.enabled {
		return
	}
	go func() {
		c.mutex.Lock()
		if op, ok := c.operations[key]; ok && op.eventReceived {
			fmt.Fprintf(os.Stderr, "VERIFIER:expectUpdate: package operation for %s already exists (event received first)\n", key)
			if op.resourceVersion <= resourceVersion {
				fmt.Fprintf(os.Stderr, "VERIFIER:expectUpdate: update verified – deleting operation\n")
				delete(c.operations, key)
			}
		} else {
			if ok && !op.eventReceived {
				fmt.Fprintf(os.Stderr, "VERIFIER:expectUpdate: package operation for %s already exists (event NOT received yet)\n", key)
			}
			fmt.Fprintf(os.Stderr, "VERIFIER:expectUpdate: registering package operation for %s with min resource version %s\n",
				key, resourceVersion)
			c.operations[key] = operation{
				timestamp:       time.Now(),
				resourceVersion: resourceVersion,
			}
		}
		c.mutex.Unlock()
	}()
}

func (c *connectionVerifier) expectDelete(key string) {
	if !c.enabled {
		return
	}
	// TODO
}

func (c *connectionVerifier) addReceived(key string, resourceVersion string) {
	if !c.enabled {
		return
	}
	// TODO
}

func (c *connectionVerifier) updateReceived(key string, resourceVersion string) {
	if !c.enabled {
		return
	}
	go func() {
		c.mutex.Lock()
		fmt.Fprintf(os.Stderr, "VERIFIER:updateReceived: %s %s\n", key, resourceVersion)
		if op, ok := c.operations[key]; ok {
			if !op.eventReceived && op.resourceVersion <= resourceVersion {
				delete(c.operations, key)
				fmt.Fprintf(os.Stderr, "VERIFIER:updateReceived: update verified (deleted) for %s\n", key)
			} else {
				if op.eventReceived {
					fmt.Fprintf(os.Stderr, "VERIFIER:updateReceived: update received but already exists for %s\n", key)
				} else {
					fmt.Fprintf(os.Stderr, "VERIFIER:updateReceived: update received but resource version smaller:"+
						"than expected: %s < %s\n", resourceVersion, op.resourceVersion)
				}
			}
		} else {
			fmt.Fprintf(os.Stderr, "VERIFIER:updateReceived: registering package operation for %s\n", key)
			c.operations[key] = operation{
				timestamp:       time.Now(),
				resourceVersion: resourceVersion,
				eventReceived:   true,
			}
		}
		c.mutex.Unlock()
	}()
}

func (c *connectionVerifier) deleteReceived(key string, resourceVersion string) {
	if !c.enabled {
		return
	}
	// TODO
}
