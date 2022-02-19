package distributed

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

type AccessControl struct {
	// used for any custom data
	SourcesArray map[int]interface{}
	// map is not multi-thread safe, so locks are needed
	localHost string
	port      string
	hashCon   *Consistent
	*sync.RWMutex
}

func (m *AccessControl) SetHosts(localHost string, port string) {
	m.localHost = localHost
	m.port = port
}

func (m *AccessControl) SetConsistentHash(consistent *Consistent) {
	m.hashCon = consistent
}

func (m *AccessControl) GetNewRecord(uid int) interface{} {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	data := m.SourcesArray[uid]
	return data
}

func (m *AccessControl) SetNewRecord(uid int) {
	m.RWMutex.Lock()
	m.SourcesArray[uid] = "hello jzlmall"
	m.RWMutex.Unlock()
}

func (m *AccessControl) GetDistributedRight(req *http.Request) bool {
	uid, err := req.Cookie("uid")
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
		return m.GetDataFromOtherMap(hostRequest, req)
	}

}

func (m *AccessControl) GetDataFromOtherMap(host string, request *http.Request) bool {
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return false
	}
	uidSign, err := request.Cookie("sign")
	if err != nil {
		return false
	}

	// mock http API request
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://"+host+":"+m.port+"/access", nil)
	if err != nil {
		return false
	}

	// put required cookies into the mocked request
	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	coookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	req.AddCookie(cookieUid)
	req.AddCookie(coookieSign)

	// retrieve from the response
	response, err := client.Do(req)
	if err != nil {
		return false
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}

	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}
	return false
}

func (m *AccessControl) GetDataFromMap(uid string) bool {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}
	data := m.GetNewRecord(uidInt)

	if data != nil {
		return true
	}
	return false
}
