package controllers

import (
	"encoding/json"
	"errors"
	"github.com/anaskhan96/soup"
	"github.com/yamamushi/durouter/config"
	"github.com/yamamushi/durouter/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ServerStatus struct {
	config config.Config
}




func NewServerStatus(config config.Config) (*ServerStatus){
	serverstatus := &ServerStatus{}
	serverstatus.config = config
	return serverstatus
}

func (h *ServerStatus) GetStatus(w http.ResponseWriter, r *http.Request) {
	// Tell our client to expect a json response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	statuslist, err := h.GetStatusList()
	if err != nil {
		status := &models.ServerStatus{Error:true, ErrorStatus:"Could not retrieve status"}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(status)
		return
	}

	currentTest := h.FindCurrentTest(statuslist)
	if currentTest.Error {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(currentTest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(currentTest)
}

func (h *ServerStatus) GetSchedule (w http.ResponseWriter, r *http.Request) {
	// Tell our client to expect a json response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	statuslist, err := h.GetStatusList()
	if err != nil {
		status := &models.ServerStatus{Error:true, ErrorStatus:"Could not retrieve status"}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(status)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(statuslist)
}


func (h *ServerStatus) GetStatusList() (statuslist []models.ServerStatus, err error) {

	resp, err := soup.Get("https://www.dualthegame.com/en/server-status/") // Append page=1000 so we get the last page
	if err != nil {
		//fmt.Println("Could not retreive page: " + record.ForumProfile)
		return statuslist, err
	}

	doc := soup.HTMLParse(resp)

	statusTableBox := doc.FindAll("div", "class", "table-responsive")
	//fmt.Println("Tables: " + strconv.Itoa(len(statusTableBox)))

	if len(statusTableBox) > 0 {

		statusTableRows := statusTableBox[0].Find("table", "class", "table").FindAll("tr")

		//fmt.Println("Rows: " + strconv.Itoa(len(statusTableRows)))

		if len(statusTableRows) > 1 {

			for rowNum, row := range statusTableRows {
				// We always skip the first row, which is a description row
				var status models.ServerStatus
				if rowNum > 0 {
					rowColumns := row.FindAll("td")
					//fmt.Println("Row "+strconv.Itoa(rowNum)+ " Column Count: " + strconv.Itoa(len(rowColumns)))

					for colNumber, column := range rowColumns {
						//fmt.Println(strings.TrimSpace(column.Text()))
						if colNumber == 0 {
							classAttrs := column.Attrs()["class"]
							colorString := strings.Trim(strings.Split(classAttrs, " ")[0], "is-")

							if colorString == "green" {
								status.StatusColor = 6932560
							}
							if colorString == "gray" {
								status.StatusColor = 14013909
							}
							if colorString == "yellow" {
								status.StatusColor = 16380271
							}
							if colorString == "orange" {
								status.StatusColor = 16743941
							}
							if colorString == "red" {
								status.StatusColor = 14417920
							}
							status.Status = strings.TrimSpace(column.Text())
						}
						if colNumber == 1 {
							status.TestType = strings.TrimSpace(column.Text())
						}
						if colNumber == 2 {
							if len(column.Text()) > 0 {
								status.Access = strings.TrimSpace(column.Text())
							} else {
								afield := column.Find("a")
								status.Access = afield.Text()
							}
						}
						if colNumber == 3 {
							//fmt.Println(strings.TrimSpace(column.Text()))
							status.StartDate, err  = time.Parse("January 02, 2006 - 15:04 MST", strings.TrimSpace(column.Text()))
							if err != nil {
								//fmt.Println(err.Error())
								return statuslist, errors.New("Could not parse start date for row " + strconv.Itoa(rowNum))
							}
						}
						if colNumber == 4 {
							status.EndDate, err  = time.Parse("January 02, 2006 - 15:04 MST", strings.TrimSpace(column.Text()))
							if err != nil {
								return statuslist, errors.New("Could not parse end date for row " + strconv.Itoa(rowNum))
							}
						}
						if colNumber == 5 {
							timeFields := strings.Split(strings.TrimSpace(column.Text()), ":")
							if len(timeFields) < 2 {
								return statuslist, errors.New("Could not parse duration for row " + strconv.Itoa(rowNum))
							}
							hoursparsed, err := strconv.Atoi(timeFields[0])
							if err != nil {
								return statuslist, errors.New("Could not parse hour duration for row " + strconv.Itoa(rowNum) + " - " + err.Error())
							}
							minutesparsed, err := strconv.Atoi(timeFields[1])
							if err != nil {
								return statuslist, errors.New("Could not parse minute duration for row " + strconv.Itoa(rowNum) + " - " + err.Error())
							}

							hours := time.Duration(time.Hour * time.Duration(hoursparsed))
							minutes := time.Duration(time.Minute * time.Duration(minutesparsed))
							duration := time.Duration(hours + minutes)
							status.Duration = duration
						}

					}
					statuslist = append(statuslist, status)
					//fmt.Println("\n")
				}
			}

		} else {
			return statuslist, errors.New("Could not parse status table correctly")
		}

	}

	return statuslist, nil
}


func (h *ServerStatus) FindCurrentTest(statuslist []models.ServerStatus) (status models.ServerStatus) {

	for num, item := range statuslist {
		item.Error = false
		if strings.ToLower(item.Status) == "live" {
			return item
		}

		if time.Now().Before(item.StartDate) {
			if num == 0 {
				return item
			}
			if time.Now().After(statuslist[num-1].EndDate) {
				return item
			}
		}
	}

	status.Error = true
	status.ErrorStatus = "Could not find current or next test"

	return status
}

