package items

import (
	"encoding/xml"	
	"strconv"
	"io"
	"net/http"
)

// Define the structure of the XML based on the format provided
type Wowhead struct {
	Item struct {
		ID     string `xml:"id,attr"`
		Icon   struct {
			DisplayId string `xml:"displayId,attr"`
		} `xml:"icon"`
	} `xml:"item"`
}

// Function to fetch and parse the XML
func GetDisplayId(itemId int) (string, error) {

	url := "http://www.wowhead.com/?item=" + strconv.Itoa(itemId) + "&xml"

	// Fetch the XML data from the URL
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse the XML
	var wowhead Wowhead
	err = xml.Unmarshal(body, &wowhead)
	if err != nil {
		return "", err
	}

	// Extract the displayId
	displayId := wowhead.Item.Icon.DisplayId
	return displayId, nil
}

// func main() {
// 	// Example URL
// 	url := "http://www.wowhead.com/?item=33214&xml"
// 	displayId, err := GetDisplayIdFromUrl(url)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}

// 	fmt.Println("Display ID:", displayId)

// }
	
