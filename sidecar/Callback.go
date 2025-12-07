package sidecar


type Callback struct {
    Context map[string]interface{}
    Chain   []func(*Callback) error
}

// CallNext exécute et retire le premier élément de la pile
func (c *Callback) CallNext() error {
    if len(c.Chain) == 0 {
        return nil // fin de la chaîne
    }
    // prendre le premier
    fn := c.Chain[0]
    // raccourcir la slice (pop)
    c.Chain = c.Chain[1:]
    // exécuter
    return fn(c)
}


func (c *Callback) Clone() *Callback {
    // copier la map
    newCtx := make(map[string]interface{}, len(c.Context))
    for k, v := range c.Context {
        newCtx[k] = v
    }

    // copier la slice (les fonctions elles-mêmes sont des pointeurs immuables)
    newChain := make([]func(*Callback) error, len(c.Chain))
    copy(newChain, c.Chain)

    return &Callback{
        Context: newCtx,
        Chain:   newChain,
    }
}
