//author SkyRocknRoll

package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/redis.v3"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	ListenEndpoint := flag.String("listen-enpoint", "0.0.0.0:9736", "Enpoint to listen ex: 0.0.0.0:9736")
	r := mux.NewRouter()
	r.HandleFunc("/redis/{endpoint}", RedisStatus)
	r.HandleFunc("/sentinel/{endpoint}/{master}", RedisSentinelStatus)
	r.HandleFunc("/cluster/{nodeAddresses}", RedisClusterStatus)

	fmt.Printf("Server listening on %s", *ListenEndpoint)
	err := http.ListenAndServe(*ListenEndpoint, r)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Unable to start the server!!!")
	}

}

// Creates New Redis Client
func GetNewClient(Endpoint string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     Endpoint,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return client
}

//Provides new redis client using sentinel
func GetNewFailoverClient(MasterName string, SentinelAddrs []string) *redis.Client {
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    MasterName,
		SentinelAddrs: SentinelAddrs,
		Password:      "", // no password set
		DB:            0,  // use default DB
	})

	return client
}

func GetNewClusterClient(NodeAddresses []string) *redis.ClusterClient{
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: NodeAddresses,
		Password: "", // no password set
	})
	return client
}

func RedisStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	client := GetNewClient(vars["endpoint"])
	defer client.Close() // Close the client.
	data, err := client.Set("__||__", "Are You Up?", -1).Result()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		fmt.Println(err.Error())

	} else {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, data)

	}
}

func RedisSentinelStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	SentinelAddrs := strings.Split(vars["endpoint"], ",")
	MasterName := vars["master"]
	client := GetNewFailoverClient(MasterName, SentinelAddrs)
	defer client.Close() // Close the client.
	data, err := client.Set("__||__", "Are You Up?", -1).Result()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		//fmt.Println(err.Error())

	} else {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, data)

	}
}

func RedisClusterStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	NodeAddresses := strings.Split(vars["nodeAddresses"], ",")
	client := GetNewClusterClient(NodeAddresses)
	defer client.Close() // Close the client.
	client.Ping()
	data, err := client.Set("__||__", "Are You Up?", -1).Result()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		fmt.Println(err.Error())

	} else {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, data)

	}
}