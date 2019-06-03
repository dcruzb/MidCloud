package common

import (
	"github.com/dcbCIn/MidCloud/lib"
	"sort"
	"strings"
	"time"
)

type CloudService struct {
	Aor    NamingRecord
	Price  float64
	Status bool
	Rank   int
}

// Calls each of the cloud services using a generic proxy to get data to make the Rank based on a predefined sort method
func (cs *CloudService) RefreshPrice() {
	// Todo implement call to each of the cloud services using a generic proxy to get current price
}

// Calls each of the cloud services using a generic proxy to get status of them to make the Rank based on a predefined sort method
func (cs *CloudService) RefreshStatus() {
	// Todo implement call to each of the cloud services using a generic proxy to get status
}

// Type used to sort []CloudService by Price and Availability. Implements sort.Interface
type SortByPriceAndAvailability []CloudService

// Len is the number of elements in the collection.
func (s SortByPriceAndAvailability) Len() int {
	return len(s)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (s SortByPriceAndAvailability) Less(i, j int) bool {
	return (s[i].Status && s[i].Price < s[j].Price) || !s[j].Status
}

// Swap swaps the elements with indexes i and j.
func (s SortByPriceAndAvailability) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type Monitor struct {
	cloudServices []CloudService
	//nameServerIp          string
	//nameServerPort        int
	lookup                ILookup
	cloudFunctionName     string
	cloudFunctionsPattern string
}

// Starts monitoring the clouds available and caches the rank of their prices
// cloudFunctionName stands for the name that will be binded in the nameServer
//
//	monitor := common.Monitor{}
//	go monitor.Start(shared.NAME_SERVER_IP, shared.NAME_SERVER_PORT, "cloudFunctions", "CloudFunctions")
func (mon *Monitor) Start( /*nameServerIp string, nameServerPort int,*/ lookupProxy ILookup, cloudFunctionName string, cloudFunctionsPattern string) {
	//mon.nameServerIp = nameServerIp
	//mon.nameServerPort = nameServerPort
	mon.lookup = lookupProxy
	mon.cloudFunctionName = cloudFunctionName
	mon.cloudFunctionsPattern = cloudFunctionsPattern
	for {
		mon.RefreshRank()
		time.Sleep(1 * time.Minute)
	}
	// Todo monitor não está fechando e reabrindo a conexão para o lookup, pois é acessado diversas vezes depois (a cada intervalo de tempo pré-definido). Isso ocasiona a quebra do sistema de monitoramento caso o servidor de nomes seja reiniciado.
}

// Get the list of cloud services based on name server list of binded servers
func (mon Monitor) RefreshCloudServices() {
	//lp := dist.NewLookupProxy(mon.nameServerIp, mon.nameServerPort)
	services, err := mon.lookup.List()
	if err != nil {
		lib.PrintlnError("Error at lookup. Error:", err)
	}
	//err = mon.lookup.Close()
	//if err != nil {
	//	lib.PrintlnError("Error at closing lookup. Error:", err)
	//}

	for _, service := range services {
		// If the service registred in NameServer is a CloudFunctions server
		if strings.Contains(service.ServiceName, mon.cloudFunctionsPattern) {
			found := false
			for _, cloudService := range mon.cloudServices {
				if cloudService.Aor.ServiceName == service.ServiceName {
					found = true
				}
			}
			if !found {
				newCloudService := CloudService{}
				newCloudService.Aor = service
				mon.cloudServices = append(mon.cloudServices, newCloudService)
			}
		}
	}
}

// Uses the chosen sort method to generate a Rank of the best cost-efective cloud servers
func (mon *Monitor) RefreshRank() {
	mon.RefreshCloudServices()
	for _, service := range mon.cloudServices {
		service.RefreshPrice()
		service.RefreshStatus()
	}

	sort.Sort(SortByPriceAndAvailability(mon.cloudServices))
}

// Binds the best cost-benefit cloud server to de name server
func (mon *Monitor) Bind() {
	//lp := dist.NewLookupProxy(mon.nameServerIp, mon.nameServerPort)
	err := mon.lookup.Bind(mon.cloudFunctionName, mon.cloudServices[0].Aor.ClientProxy)
	lib.FailOnError(err, "Error at lookup.")
	//err = mon.lookup.Close()
	//lib.FailOnError(err, "Error at closing lookup")
}
