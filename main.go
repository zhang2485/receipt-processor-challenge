package main

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"unicode"

	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Data Models
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type DataStore struct {
	dto map[string]Receipt
}

// Data storage functionality
func (d *DataStore) StoreReceipt(receipt Receipt, id string) {
	if d.dto == nil {
		d.dto = make(map[string]Receipt)
	}
	d.dto[id] = receipt
}

func (d *DataStore) GetReceipt(id string) Receipt {
	if d.dto == nil {
		d.dto = make(map[string]Receipt)
	}
	if r, ok := d.dto[id]; ok {
		return r
	}
	return Receipt{}
}

// Data Model functionality
func (i *Item) calculatePoints() int {
	if len(strings.TrimSpace(i.ShortDescription))%3 == 0 {
		price, _ := strconv.ParseFloat(i.Price, 64)
		return int(math.Ceil(price * 0.2))
	}
	return 0
}

func (i *Item) validateData() bool {
	re := regexp.MustCompile("^[\\w\\s\\-]+$")
	if !re.Match([]byte(i.ShortDescription)) {
		return false
	}
	re = regexp.MustCompile("^\\d+\\.\\d{2}$")
	return re.Match([]byte(i.Price))
}

func (r *Receipt) validateData() bool {
	re := regexp.MustCompile("^[\\w\\s\\-&]+$")
	if !re.Match([]byte(r.Retailer)) {
		return false
	}
	re = regexp.MustCompile("^\\d+\\.\\d{2}$")
	if !re.Match([]byte(r.Total)) {
		return false
	}
	layout := "2006-01-02"
	if _, err := time.Parse(layout, r.PurchaseDate); err != nil {
		return false
	}
	if len(r.PurchaseTime) == 4 {
		r.PurchaseTime = "0" + r.PurchaseTime
	}
	if matched, _ := regexp.MatchString(`^([01]\d|2[0-3]):([0-5]\d)$`, r.PurchaseTime); !matched {
		return false
	}
	for _, i := range r.Items {
		if !i.validateData() {
			return false
		}
	}
	return len(r.Items) > 0
}

func (r *Receipt) calculatePoints() int {
	score := 0
	for _, c := range r.Retailer {
		if unicode.IsDigit(c) || unicode.IsLetter(c) {
			score += 1
		}
	}
	receiptTotal, _ := strconv.ParseFloat(r.Total, 64)
	if receiptTotal-float64(int(receiptTotal)) == 0 {
		score += 50
	}
	if int(receiptTotal*100)%25 == 0 {
		score += 25
	}
	for _, item := range r.Items {
		score += item.calculatePoints()
	}
	score += (len(r.Items) / 2) * 5
	date, _ := time.Parse("2006-01-02", r.PurchaseDate)
	if date.Day()%2 == 1 {
		score += 6
	}
	purchaseTime := r.PurchaseTime
	if len(purchaseTime) == 4 {
		purchaseTime = "0" + purchaseTime
	}
	if purchaseTime > "14:00" && purchaseTime < "16:00" {
		score += 10
	}
	return score
}

// Endpoint Functionality
var dataStore = DataStore{}

func getPoints(c *gin.Context) {
	id := c.Param("id")
	receipt := dataStore.GetReceipt(id)
	if len(receipt.Retailer) == 0 {
		c.IndentedJSON(http.StatusNotFound, "#/components/responses/NotFound")
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"points": receipt.calculatePoints()})
}

func processReceipt(c *gin.Context) {
	id := uuid.New()
	var newReceipt Receipt
	if err := c.BindJSON(&newReceipt); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "#/components/responses/BadRequest")
		return
	}
	if !newReceipt.validateData() {
		c.IndentedJSON(http.StatusBadRequest, "#/components/responses/BadRequest")
		return
	}
	dataStore.StoreReceipt(newReceipt, id.String())
	c.IndentedJSON(http.StatusCreated, gin.H{"id": id.String()})
}

func main() {
	router := gin.Default()
	router.POST("/receipts/process", processReceipt)
	router.GET("/receipts/:id/points", getPoints)

	router.Run("localhost:8080")
}
