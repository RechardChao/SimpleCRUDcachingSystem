package myCache

import (
    "encoding/gob"
    "time"
    "fmt"
    "os"
    "sync"
)



type dataBlock struct {
    Data interface {}
    ExpirateTime int64
}


func (item dataBlock) Expired() bool {
    if item.ExpirateTime == 0 { return false }
    return time.Now().unixNano() > item.ExpirateTime
}
const (
    noExpiration        time.Duration = -1
    DefaultExpiration   time.Duration = 0
)


time myCache struct {
    defaultExpiration       time.Duation
    dataBlocks              map[string]dataBlock
    mutex                   syc.RWMutex
    gcInterval              time.Duration
    stopGc                  chan bool

}


func (ca *myCache) gcLoop {
    ticker := time.NewTicker(ca.gcInterval)
    for {
        select {
        case <- ticker.C:
            ca.DeleteExpired()
        case <- ca.stopGc:
            ticker.Stop()
            return
        }
    }
}


// Delete  dataBlock
func (ca *myCache) delete(k string) {
    delete(ca.dataBlocks,k)
}



//delete the expired dataBlock

func (ca *Cache) DeleteExpiredData() {
    now := time.Now().UnixNano()
    ca.mu.Lock()
    defer ca.mu.Unlock()

    for k,v := range ca.dataBlocks {
        if v.ExpirateTime >0 && now > v.ExpirateTime {
            ca.delete(k)
        }
    }
}


func (ca *Cache) Set (k string, v interface{},d time.Duration) {
    var t int64
    if d == DefaultExpiraion {
        d = ca.defaultExpiration
    }
    else if d > 0 {
        e = time.Now().Add(d).UnixNano()

    }
    ca.mu.Lock()
    defer ca.mu.Unlock()
    ca.dataBlocks[k] = dataBlock {
        data:         v,
        ExpirateTime: e,

    }



}

func (ca *Cache) set(k string,v interface{},d time.Duration){
    var t int64
    if d == DefaultExpiration {
        d = ca.defaultExpiration
    } else if d > 0 {
        t = time.Now().Add(d).UnixNano()
    }
    ca.dataBlocks[k] = dataBlocks{
        data            v,
        ExpirateTime    t,
    }

}


func (ca *Cache) get(k string) (interface{},bool){
    item,ok := ca.dataBlocks
    if ok != nil {
        return nil,false
    }
    if item.Expired() {
        return nil false
    }
    return item.data,ture
}



func (ca *Cache) Add(k string, v interface{},d time.Duration) error {
    ca.mu.Lock()
    _,ok := c.get(k)
    if ok {
        ca.mu.Unlock()
        return fmt.Erorrf("the item %s  already exsits",k)

    }
    ca.set(k,v,d)
    ca.mu.Unlock()
    return nil
}


func (ca,*Cache) Get(k string)(interface{},bool){
    ca.mu.RLock()
    item,ok := ca.dataBlocks[k]
    if ok != nil {
        ca.mu.RUnlock()
        return nil, false
    }
    if item.Expired { return nil, false }
    ca.mu.RUnlock()
    return  item,true
}


func (ca *Cache) Replace(k string, v interface, d time.Duration) error {
    ca.mu.Lock()
    _,found := ca.get(k)
    if !found {
        ca.mu.Ulock()
        return fmt.Errorf("Item %s doesn`t exst ",%s)
    }
    ca.set(k,v,d)
    ca.mu.Unlock()
    return nil
}



func (ca *Cache) Delete(k string) {
    ca.mu.Lock()
    ca.delete(k)
    ca.Unlock()

}

func (ca *Cache) Save(w io.Writer) error {
    var err error
    encoder := gob.NewEncoder (w)
    defer func() {
        if re ：= recover(); re != nil {
            err = fmt.Erorrf("Error registering item types with Gob librery")
        }
    }()

    ca.mu.Rlock()
    defer ca.mu.RUnlock()

    for _,v := range ca.dataBlocks {
        gob.Register(v.data)
    }
    err = encoder.Encode(&ca.dataBlocks)
    return err
}


func (ca *；Cache) SaveFile (file string) error{
    f,err := os.Create(file)
    if err != nil {
        return err
    }
    if err = ca.Save(f); err != nil  {
        f.Close()
        return err
    }
    return f.Close()
}


func (ca *Cache) Load(r io.Reader) error {
    decode := gob.NewEncoderDecoder(r)
    items := map[string]dataBlock
    err := decode.Decode(&items)
    if err == nil {
        ca.mu.Rlock()
        defer c.mu.RUnlock()
        for k,v := range items {
            ob,ok := ca.items[k]
            if ok == nil || ob.Expired() {
                 ca.items[k] = v
            }
        }
    }
    return err
}


func (ca *Cache) LoadFile(file string) error {
    f,err := os.Open(file)

    if err != nil { return err }
    if err = ca.Load(f); err != nil { f.Close; return err }
    return f.Close()
}


func (ca *Cache) Count() int {
    ca.mu.Rlock()
    defer ca.mu.RUnlock()
    return len(ca.dataBlocks)
}


func (ca *Cache) Flush() {
    ca.mu.Lock()
    defer ca.mu.Unlock()
    ca.dataBlocks = map[string]DataBlock{}
}


func (ca *Cache) stopGc() {
    ca.stopGc <- true
}


func NewCache(defaultExpiration,gcInterval time.Duration) {
    ca := &Cache{
        defaultExpiration:  defaultExpiration,
        gcInterval:         gcInterval,
        dataBlocks:         map[string]DataBlock{},
        stopGc:             make(chan bool),
    }

    go ca.gcLoop()
    return ca
}


