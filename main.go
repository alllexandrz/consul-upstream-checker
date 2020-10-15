
package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	logFile, err := os.OpenFile("log.log", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	consulAddress := os.Getenv("CONSUL")
	dcName := os.Getenv("DC")
	log.Printf("Will be working with %s %s",consulAddress,dcName)
	client, err := api.NewClient(&api.Config{Address: consulAddress})
	if err != nil {
		panic(err)
	}

	catalog := client.Catalog()
	q := api.QueryOptions{
		WaitIndex: 0,
		WaitTime:  time.Second * 10,
	}
		services, meta, err := catalog.Services(&q)
		log.Printf("Willbe processed %d services",len(services))
		if err != nil {
			log.Fatal("Failed to get service catalog from consul agent: ", err)
		}
		for svcName := range services {
			svcCatalog, _, err := catalog.Service(svcName, "", nil)
			if err != nil {
				log.Fatal("Failed to get service entry from consul agent: ", err)
			}
			if len(svcCatalog) == 0 {
				continue
			}
			svc := svcCatalog[0]
			//	Send http request
			timeout := time.Duration(3 * time.Second)
			client := http.Client{
				Timeout: timeout,
			}
			serviceUrl := fmt.Sprintf("http://%s:%d", svc.Address,svc.ServicePort)
			resp, err := client.Get(serviceUrl)
			if err != nil {
				log.Printf("ERROR: %s:%d (%s) %v",svc.Address,svc.ServicePort,svc.ServiceName,err)
			}else{
				// Print the HTTP Status Code and Status Name
				if resp.StatusCode != 200 {
					log.Printf("HTTP %s:%d (%s) Response Status: [%d] %s",svc.Address,svc.ServicePort,svc.ServiceName, resp.StatusCode, http.StatusText(resp.StatusCode))
				}
			}

			q.WaitIndex = meta.LastIndex

	}



	router := mux.NewRouter()
	router.HandleFunc("/",handleRequestAndRedirect )
	http.Handle("/",router)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":5000", nil)

}
func handleRequestAndRedirect(w http.ResponseWriter, r *http.Request) {

	log.Printf("finished.")
}
