package service

import (
	"fmt"
	"net/http"
	"svclookup/xutil"
	"sync"
)

func LookupSvc(URLs []string) []xutil.RespItem {
	var client = GetClient(); 
	var respArr = []xutil.RespItem{}

	for i := range URLs {
		// fmt.Printf("Index: %d \n", i)
		var site, sCode, desc1 = xutil.ProcessReqMV(URLs[i], client)
		respArr = append(respArr, xutil.RespItem{
			Site: site,
			StatusCode: sCode,
			Description: desc1,
		})
	}

	return respArr
}

func reqWorker (res chan<- xutil.RespItem, url string, wg *sync.WaitGroup) {
	defer wg.Done()

	var client = GetClient();
	var respItem = xutil.ProcessReq(url, client)
	if respItem.StatusCode != http.StatusOK {
		respDet := xutil.RespItem(respItem)
		res <- respDet
	}
	fmt.Println("REQ :: ", url);
}

func LookupSvcAlt(URLs []string) []xutil.RespItem {
	ch := make(chan xutil.RespItem, len(URLs))
	var wg sync.WaitGroup;
	
	for i := range URLs {
		wg.Add(1)
		go reqWorker(ch, URLs[i], &wg)
	}
	
	var respArr []xutil.RespItem
	// var respArr []xutil.RespItem{}

	var wg2 sync.WaitGroup
    wg2.Add(1)
	go func() {
		defer wg2.Done()

		for range URLs {
			response := <-ch
			if response.StatusCode != http.StatusOK {
				respArr = append(respArr, response)
			}
		}
    }()

	wg.Wait()
    close(ch)
	wg2.Wait()
	
	return respArr
}

/*
func LookupSvcAlt(URLs []string) []xutil.RespItem {
	ch := make(chan xutil.RespItem, len(URLs))
	var wg sync.WaitGroup;
	
	for i := range URLs {
		wg.Add(1)

		go func () {
			var client = GetClient();
			var respItem = xutil.ProcessReq(URLs[i], client)
			if respItem.StatusCode != http.StatusOK {
				respDet := xutil.RespItem(respItem)
				ch <- respDet
			}
			wg.Done()
		}()
	}
	wg.Wait()
	
	var respArr []xutil.RespItem

	go func() {
		for range URLs {
			response := <-ch
			if response.StatusCode != http.StatusOK {
				respArr = append(respArr, response)
			}
		}
    }()

	wg.Wait()
    close(ch)
	
	fmt.Println(respArr);
	return respArr
}*/

func GetClient() http.Client {
	transport := http.Transport{
		Dial:  xutil.ReqTimeout,
	}

	client := http.Client{
		Transport: &transport,
	}
	return client
}
