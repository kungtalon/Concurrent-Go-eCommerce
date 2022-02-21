package distributed

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// users are not allowed to buy again within 20 seconds
const interval = 20

type AccessControl struct {
	// used for any custom data
	SourcesArray map[int]time.Time
	// map is not multi-thread safe, so locks are needed
	localHost string
	port      string
	hashCon   *Consistent
	BlackList *BlackList
	sync.RWMutex
}

type BlackList struct {
	Lookup map[int]bool
	sync.RWMutex
}

func (bl *BlackList) GetBlackListByID(uid int) bool {
	bl.RLock()
	defer bl.RUnlock()
	return bl.Lookup[uid]
}

func (bl *BlackList) SetBlackListByID(uid int) bool {
	bl.Lock()
	defer bl.Unlock()
	bl.Lookup[uid] = true
	return true
}

func (m *AccessControl) SetHosts(localHost string, port string) {
	m.localHost = localHost
	m.port = port
}

func (m *AccessControl) SetConsistentHash(consistent *Consistent) {
	m.hashCon = consistent
}

func (m *AccessControl) GetNewRecord(uid int) time.Time {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	data := m.SourcesArray[uid]
	log.Println("Latest access for user " + strconv.Itoa(uid) + " at time " + data.Format("2020-01-01 15:02:30"))
	return data
}

func (m *AccessControl) SetNewRecord(uid int) {
	m.RWMutex.Lock()
	m.SourcesArray[uid] = time.Now()
	m.RWMutex.Unlock()
}

func (m *AccessControl) GetDistributedRight(req *http.Request, productIdStr string) bool {
	uid, err := req.Cookie("uid")
	log.Println("Cur uid: " + uid.Value)
	if err != nil {
		return false
	}

	// get the closest server on hash ring
	hostRequest, err := m.hashCon.Get(uid.Value)
	if err != nil {
		return false
	}

	// check whether target is local
	if hostRequest == m.localHost {
		return m.GetDataFromMap(uid.Value)
	} else {
		return m.GetDataFromOtherMap(hostRequest, req, productIdStr)
	}

}

func (m *AccessControl) GetUrl(hostUrl string, request *http.Request) (response *http.Response, body []byte, err error) {
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return
	}
	uidSign, err := request.Cookie("sign")
	if err != nil {
		return
	}

	// mock http API request
	// TODO: use gRPC instead of HTTP for communication
	client := &http.Client{}
	req, err := http.NewRequest("GET", hostUrl, nil)
	if err != nil {
		return
	}

	// put required cookies into the mocked request
	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	coookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	req.AddCookie(cookieUid)
	req.AddCookie(coookieSign)

	// retrieve from the response
	response, err = client.Do(req)
	defer response.Body.Close()
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(response.Body)
	return
}

func (m *AccessControl) GetDataFromMap(uid string) bool {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}

	if m.BlackList.GetBlackListByID(uidInt) {
		return false
	}

	dataRecord := m.GetNewRecord(uidInt)
	if !dataRecord.IsZero() {
		if dataRecord.Add(time.Duration(interval) * time.Second).After(time.Now()) {
			log.Println("User " + uid + " operated too frequently.")
			return false
		}
	}
	m.SetNewRecord(uidInt)
	return true
}

func (m *AccessControl) GetDataFromOtherMap(host string, request *http.Request, productIdStr string) bool {
	hostUrl := "http://" + host + ":" + m.port + "/check?productID=" + productIdStr
	log.Println("Sending check request to " + host)
	response, body, err := m.GetUrl(hostUrl, request)
	if err != nil {
		return false
	}

	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		}
	}
	return false
}
