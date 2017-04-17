package cluster

import (
	"strings"
	"sync"
	"time"
)

// AutoPurgeInterval is the maximum amount of time to wait before automatically purging the cluster of stale members
var AutoPurgeInterval = 24 * time.Hour

// AutoPurgeAge is the maximum age allowed for members' last update when automatically purging.
var AutoPurgeAge = 12 * time.Hour

// Cluster describes the set of ari proxies in a system.  The list is indexed by a hash of the asterisk ID and the ARI application and indicates the time of last contact.
type Cluster struct {
	lastPurge time.Time

	members map[string]time.Time

	mu sync.Mutex
}

// New returns a new Cluster
func New() *Cluster {
	return &Cluster{
		members: make(map[string]time.Time),
	}
}

// hash returns the key for a given proxy instance
func hash(id, app string) string {
	return id + "|" + app
}

// dehash returns the Asterisk ID and ARI application represented by the given key
func dehash(key string) (id string, app string) {
	pieces := strings.Split(key, "|")
	if len(pieces) < 2 {
		return
	}
	return pieces[0], pieces[1]
}

// Member describes the state of a Member of an application cluster
type Member struct {
	// ID is the unique identifier for the Asterisk node
	ID string

	// App indicates the ARI application of this proxy
	App string

	// LastActive is the timestamp of the last occurrence of this node
	LastActive time.Time
}

// All returns a list of all cluster members whose LastActive time is no older thatn the given maxAge.
func (c *Cluster) All(maxAge time.Duration) (list []Member) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.members {
		if maxAge == 0 || time.Since(v) < maxAge {
			id, app := dehash(k)
			list = append(list, Member{
				ID:         id,
				App:        app,
				LastActive: v,
			})
		}
	}
	return
}

// App returns a list of all cluster members for the given ARI Application whose LastActive time is no older than the given maxAge.
func (c *Cluster) App(app string, maxAge time.Duration) (list []Member) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.members {
		i, a := dehash(k)
		if app == a && (maxAge == 0 || time.Since(v) < maxAge) {
			list = append(list, Member{
				ID:         i,
				App:        a,
				LastActive: v,
			})
		}
	}
	return
}

// Matching returns a list of all cluster members for whom the given proxy Metadata matches
func (c *Cluster) Matching(id, app string, maxAge time.Duration) (list []Member) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.members {
		if time.Since(v) > maxAge {
			continue
		}

		i, a := dehash(k)
		if id != "" && id != i {
			continue
		}
		if app != "" && app != a {
			continue
		}
		list = append(list, Member{
			ID:         i,
			App:        a,
			LastActive: v,
		})
	}
	return
}

// Update adds (or updates) a proxy to/in the cluster
func (c *Cluster) Update(id, app string) {
	c.mu.Lock()
	c.members[hash(id, app)] = time.Now()
	c.mu.Unlock()

	// See if it is time to auto-purge
	if time.Since(c.lastPurge) > AutoPurgeInterval {
		c.Purge(AutoPurgeAge)
	}
}

// Purge removes any proxies in the cluster which are older than the given maxAge.
func (c *Cluster) Purge(maxAge time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lastPurge = time.Now()

	var removalKeys []string

	for k, v := range c.members {
		if maxAge == 0 || time.Since(v) > maxAge {
			removalKeys = append(removalKeys, k)
		}
	}

	for _, key := range removalKeys {
		delete(c.members, key)
	}
}
