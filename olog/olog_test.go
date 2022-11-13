package olog

import (
	"testing"
)

func Test(t *testing.T) {
    var log = &Olog{
        Level: LEVEL_INFO,
    }
    log.Update()
    
}