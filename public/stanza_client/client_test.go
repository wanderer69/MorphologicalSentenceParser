package natashaclient

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	ctx := context.Background()
	nc := NewStanzaClient()
	require.NoError(t, nc.Init())
	time.Sleep(time.Duration(2) * time.Second)
	phrase1 := "грустный робот медленно едет в маленькую деревню на зеленом фургоне из леса"
	result, err := nc.ParsePhrase(ctx, phrase1)
	require.NoError(t, err)
	var data []interface{}
	require.NoError(t, json.Unmarshal([]byte(result), &data))
	fmt.Printf("%#v\r\n", data)
	require.NoError(t, nc.Close())
}
