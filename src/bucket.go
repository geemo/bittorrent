package dht

import "time"

// Bucket struct
type Bucket struct {
	lastChange time.Time
}

func (b *Bucket) update()  {}
func (b *Bucket) save()    {}
func (b *Bucket) insert()  {}
func (b *Bucket) delete()  {}
func (b *Bucket) find()    {}
func (b *Bucket) closest() {}
