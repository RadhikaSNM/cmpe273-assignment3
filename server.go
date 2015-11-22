/*
Radhika SNM
009426196
*/
package main

import (
"fmt"
"gopkg.in/mgo.v2"
"gopkg.in/mgo.v2/bson"
"regexp"
"strings"
"github.com/julienschmidt/httprouter"
"net/http"
"encoding/json"
"errors"
"io/ioutil"
"net/url"
"strconv"
"bytes"
)

var dbURL string ="mongodb://DBTestUser:qwerty@ds041871.mongolab.com:41871/cmpe273database"
var ErrNotFound string= "not found"
var locationDBName string="cmpe273database"
var locationCollectionName string="locations"
var tripCollectionName string="trips"
var authToken string ="Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicHJvZmlsZSIsInJlcXVlc3RfcmVjZWlwdCIsInJlcXVlc3QiLCJoaXN0b3J5X2xpdGUiXSwic3ViIjoiZDQ0MmU3MmUtYTM4My00YzQ4LTk4NjMtNThmNmMwNTNkMGY0IiwiaXNzIjoidWJlci11czEiLCJqdGkiOiJjNTMzNTE3ZS1iZjk0LTRmYzMtOGY3Yy1lZTFlMTQzZDY5ODciLCJleHAiOjE0NTA1OTg5NjAsImlhdCI6MTQ0ODAwNjk2MCwidWFjdCI6ImdQYkxCa0lnMU9QdmZQeFNxTWhnR1o1QnZRZFFRbyIsIm5iZiI6MTQ0ODAwNjg3MCwiYXVkIjoiYy1pc09fT3UzUnFmdWk3b0NhN3BQR0twNkRaUUJqQjQifQ.A77r5DnmDESKA2Uif6n3smM54TwN5jc0L7D9Yf0Dxrt2DO6ovWE8izfe3nma2PmYri33yvGhJ5PX-3Y9c6RUC5A7sAjzkiX9BR4SOmQsIdkSHwkw6ZwzdODdUbXYQsx-BcGlLvhg9Q4AE_dquQ4OGEOZP-vgN5ZZD84dv-uAVk_zaDaCjAPCyef5V6qJzdSFUocnHzqwOE3AW0i7ZYjnEUM1IPH9o7ezO2EKckZe9jgMDb4xeONqCxvuDi2wnYY4KbFm2kjh_Lg9HxTP-3MY8JDmbJ4g_FNol_C9WxgFRYSvgbhF9_iixUo9ynuhJMpfg0d0f-CS1JJuSiuZcmz1Jg"

//location: insert into db struct
type LocationInsert struct {
 Name string `json:"name"`
 Address string `json:"address"`
 City string `json:"city"`
 State string `json:"state"`
 Zip string `json:"zip"`

}

//location: db respose struct
type LocationDBResponse struct {
    Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
    Name string `json:"name"`
    Address string `json:"address"`
    City string `json:"city"`
    State string `json:"state"`
    Zip string`json:"zip"`
    Coordinate Coordinates `json:"coordinate"`

    
}

//Coordinate struct
type Coordinates struct{
    Lat float64 `json:"lat"`
    Lng float64 `json:"lng"`
} 

//Response to be sent to the user - struct
type LocationInsertResponse struct {
    Id string `json:"id"`
    Name string `json:"name"`
    Address string `json:"address"`
    City string `json:"city"`
    State string `json:"state"`
    Zip string `json:"zip"`
    Coordinate Coordinates `json:"coordinate"`

}

//Google api result struct
type googleLocationResults struct{
    Results []struct{
        Geometry struct{
           Location struct{
            Lat float64 `json: lat`
            Lng float64  `json: lng`         
            }  `json: location`

            }  `json: geometry`
            } `json: results`
        }


//Error json struct
        type errorResponse struct{
            ErrorMessage string `json: errorMessage`
        }






        func createLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
            var errDecode error
            locationDetail:=LocationInsert{}

    //decode the sent json
            errDecode=json.NewDecoder(req.Body).Decode(&locationDetail)
            if errDecode!=nil{
                fmt.Println(errDecode.Error())
                msg:="Json sent was Empty/Incorrect .Error: "
                errorCheck(msg,rw)
                return
            }

            name:=locationDetail.Name
            address:=locationDetail.Address
            city:=locationDetail.City
            state:=locationDetail.State
            zip:=locationDetail.Zip
            fmt.Println("Name is :"+name,"Address: "+address)

   //Check if any of the expected fields are empty
            if(name==""||address==""||city==""||state==""||zip==""){
             msg1:="One or more of the fields in the sent json is missing "
             errorCheck(msg1,rw)
             return
         }


    //************************************
    //Calling the google api

         fullAdd:=address+","+city+","+state+","+zip
    //lat,long,errAdd:= getLatLong("1600 Amphitheatre Parkway Mountain View CA")
         lat,long,errAdd:= getLatLong(fullAdd)       

         if errAdd!=nil{
            errorCheck(errAdd.Error(),rw) 
            return

        }else{
            fmt.Println(lat)
            fmt.Println(long) 
        }

    //*******************************************
    //calling mongo lab

        session,c, err := connectToDB(dbURL,locationDBName,locationCollectionName)
        if err!=nil{
            errorCheck("Database connection error.",rw)
            return}
            defer session.Close()

            i := bson.NewObjectId()
            fmt.Println(i)
            fmt.Println("String version")
            idString:=i.String()
            fmt.Println(idString)

       //Extracting the ID
            r, err := regexp.Compile(`"[a-z0-9]+"`)
            split1:=r.FindString(idString)
            ID:=strings.Trim(split1,"\"")
            fmt.Println(ID)

            coord:=Coordinates{lat,long}
    //Inserting into the db
            d:=LocationDBResponse{i,name,address,city,state,zip,coord}
            err=c.Insert(d);

            if err != nil {
    //log.Fatal(err)
               fmt.Println("Database insertion error: ",err.Error())
               errorCheck("Database insertion error.",rw)
               return

           }
    //creation of the json response
           respStruct:=LocationInsertResponse{ID,name,address,city,state,zip,coord}
    //marshalling into a json

           respJson, err4 := json.Marshal(respStruct)
           if err4!=nil{
            fmt.Print("Error occcured in marshalling")
        }

    //sending it in the response
        rw.Header().Set("Content-Type","application/json")
        rw.WriteHeader(http.StatusCreated)
        fmt.Fprintf(rw, "%s", respJson)

    }




    func getLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

        idString:=p.ByName("location_id")

    //Check if given id is a valid hexademical string
        err_Hex:=checkHexString(idString)
        if err_Hex!=nil{
            errorCheck("The given location is not in the correct format",rw)
            return
        }

    //Obtaining the json from DB

    //MongoLab connection
        session,c, err := connectToDB(dbURL,locationDBName,locationCollectionName)
        if err!=nil{
            errorCheck("Database connection error.",rw)
            return}
            defer session.Close()


            result := LocationDBResponse{}
            err2 := c.Find(bson.M{"_id": bson.ObjectIdHex(idString)}).One(&result)
            fmt.Println(result)
            if err2 != nil {
        //log.Fatal(err2)
                errMsg:=err2.Error()
                fmt.Println("inside get- error")
                if err2.Error()==ErrNotFound{
                    errMsg="The given location id is incorrect. Please verify."
                }    
                errorCheck(errMsg,rw)
                return
            }

            fmt.Println("Name:", result.Name)
            fmt.Println("address",result.Address)
            fmt.Println("id",result.Id)


    //Marshling values into json
     //creating the json response
            respStruct:=LocationInsertResponse{idString,result.Name,result.Address,result.City,result.State,result.Zip,result.Coordinate}
            respJson, err4 := json.Marshal(respStruct)
            if err4!=nil{
                fmt.Print("Error occcured in marshalling")     
            }

    //sending it in the response
            rw.Header().Set("Content-Type","application/json")
            rw.WriteHeader(http.StatusOK)
            fmt.Fprintf(rw, "%s", respJson)

        }



        func updateLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params){

            locationDetail:=LocationInsert{}

    //decode the sent json
            errDecode:=json.NewDecoder(req.Body).Decode(&locationDetail)
            if errDecode!=nil{
                fmt.Println(errDecode.Error())
                msg:="Json sent was Empty/Incorrect"
                errorCheck(msg,rw)
                return
            }

            address:=locationDetail.Address
            city:=locationDetail.City
            state:=locationDetail.State
            zip:=locationDetail.Zip

            fmt.Println("Address: "+address,"city "+city,"state "+state)


        //Check if any of the expected fields are empty
            if(address==""||city==""||state==""||zip==""){
             msg1:="One or more of the fields in the sent json is missing "
             errorCheck(msg1,rw)
             return
         }



    //Obtain the id to modify:
         idString:=p.ByName("location_id")

        //Check if given id is a valid hexademical string
         err_Hex:=checkHexString(idString)
         if err_Hex!=nil{
            errorCheck("The given location is not in the correct format",rw)
            return
        }

    //************************************
    //Calling the google api

        fullAdd:=address+","+city+","+state+","+zip
    //lat,long,errAdd:= getLatLong("1600 Amphitheatre Parkway Mountain View CA")
        lat,long,errAdd:= getLatLong(fullAdd)       

        if errAdd!=nil{
            errorCheck(errAdd.Error(),rw) 
            return
        //TODO work about returning an empty
        }else{
            fmt.Println(lat)
            fmt.Println(long) 
        }

    //*******************************************

//MongoLab connection
        session,c, err := connectToDB(dbURL,locationDBName,locationCollectionName)
        if err!=nil{
            errorCheck("Database connection error.",rw)
            return}
            defer session.Close()


    //Modifying the data
            id:=bson.ObjectIdHex(idString)
            errUpdate := c.UpdateId(id,bson.M{"$set": bson.M{"address": address,"city":city,"state":state,"zip":zip,"coordinate.lat":lat,"coordinate.lng":long}})

            if errUpdate!=nil{
                fmt.Println("Inside update error")
                errMsg:=errUpdate.Error()
                if errUpdate.Error()==ErrNotFound{
                    fmt.Println("Inside ErrNotFound")
                    errMsg="The given location id is incorrect. Please verify."}
                    errorCheck(errMsg,rw)
                    return
                }

    //try to find the document again for the name:
                result := LocationDBResponse{}
                err2 := c.Find(bson.M{"_id": id}).One(&result)
                fmt.Println(result)
                if err2 != nil {
                    errorCheck(err2.Error(),rw)
                }

                coord:=Coordinates{lat,long}

    //send the json response
                respStruct:=LocationInsertResponse{idString,result.Name,address,city,state,zip,coord}
    //marshalling into a json

                respJson, err4 := json.Marshal(respStruct)
                if err4!=nil{
                    fmt.Print("Error occcured in marshalling")
                }

    //sending it in the response
                rw.Header().Set("Content-Type","application/json")
                rw.WriteHeader(http.StatusCreated)
                fmt.Fprintf(rw, "%s", respJson)

            }






            func deleteLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params){

                idString:=p.ByName("location_id")


         //Check if given id is a valid hexademical string
                err_Hex:=checkHexString(idString)
                if err_Hex!=nil{
                    errorCheck("The given location is not in the correct format",rw)
                    return
                }


                id:=bson.ObjectIdHex(idString)

    //*******************************************

    //MongoLab connection
                session,c, err := connectToDB(dbURL,locationDBName,locationCollectionName)
                if err!=nil{
                    errorCheck("Database connection error.",rw)
                    return}
                    defer session.Close()



    //Delete the id:
                    errDel:=c.RemoveId(id)
                    if errDel!=nil{
                     errMsg:=errDel.Error()
                     if errDel.Error()==ErrNotFound{
                        errMsg="The given location id is incorrect. Please verify."}
                        errorCheck(errMsg,rw)
                        return
                    }

    //Set the response
                    rw.WriteHeader(http.StatusOK)
    //fmt.Fprintf(rw, "%s", respJson)

                }


                func main() {
                    fmt.Println("=========================")
                    mux := httprouter.New()
                    mux.POST("/locations",createLocation)
                    mux.GET("/locations/:location_id", getLocation)
                    mux.PUT("/locations/:location_id", updateLocation)
                    mux.DELETE("/locations/:location_id", deleteLocation)

                    //Uber API
                    mux.POST("/trips",postPlanTrip)
                    mux.GET("/trips/:trip_id",getTrip)
                    mux.PUT("/trips/:trip_id/request",requestTrip)


                    server := http.Server{
                        Addr:        "0.0.0.0:8080",
                        Handler: mux,
                    }

                    server.ListenAndServe()

                }


 //Send error json          
                func errorCheck(errMsg string,rw http.ResponseWriter){
    //Creating a errorJson
                    errorSt:=errorResponse{errMsg}

                    errorJson, err4 := json.Marshal(errorSt)
                    if err4!=nil{
                        fmt.Print("Error occcured in marshalling")
                    }
                    rw.Header().Set("Content-Type","application/json")
                    rw.WriteHeader(http.StatusBadRequest)
                    fmt.Fprintf(rw, "%s", errorJson)


                }

//Get lat longitude from google api
                func getLatLong(CombinedAddress string) (float64,float64,error){
                    addressEnc:=url.QueryEscape(CombinedAddress)
                    fmt.Println("combined link"+addressEnc)
                    link:="http://maps.google.com/maps/api/geocode/json?address="+addressEnc
                    fmt.Println(link)
                    resp, err := http.Get(link);
                    if err != nil {
                        fmt.Println(err.Error())
                        err_google:=errors.New("Google map api connection could not be established!")
                        return 0,0,err_google
                    }

                    defer resp.Body.Close()
                    body, err1 := ioutil.ReadAll(resp.Body)
                    if err1 != nil {
                        return 0,0,err1

                    }
                    var location googleLocationResults
                //Unmarshall the response into a json
                    err2:=json.Unmarshal(body,&location)
                    if err2 != nil {
                        return 0,0,err2
                    }


    //check is results is null
                    if (len(location.Results)==0){
                        err_noRes:=errors.New("The provided address is invalid. Please Check")
                        return 0,0,err_noRes 

                    }


                    lat:=location.Results[0].Geometry.Location.Lat
                    lng:=location.Results[0].Geometry.Location.Lng

                    return lat,lng,nil
                }


//Function to connect to the database
                func connectToDB(dbURL string, dbName string, collectionName string) (*mgo.Session,*mgo.Collection,error){
                  session, err := mgo.Dial(dbURL)
                  if err != nil {
                    fmt.Println("Database connection error: ",err.Error())
                    return nil,nil,err

                }
        // Optional. Switch the session to a monotonic behavior.
                session.SetMode(mgo.Monotonic, true)
                c := session.DB(dbName).C(collectionName)  
                return session,c,nil
            }

            func checkHexString(id string) error{
             stringFormat,err:=regexp.MatchString("^[A-Fa-f0-9]{24}$", id)
             if err!=nil{
                return err
            }else if(!stringFormat){
                err_Format:=errors.New("Given location id is not in a valid format.")
                return err_Format 
            }else{
                return nil
            }


        }


        //***************************************** STARTING WITH UBER STUFF ************************************

        type test struct{
            Name string `json:"name"`
        }


    type PriceEstimatesUber struct {
    Prices         []Price `json:"prices"`
}

type PutTrip struct{
    Id  string  `json:"id"`
    Status string `json:"status"`
    StartingLocation string `json:"starting_from_location_id"`
    NextLocation string `json:"next_destination_location_id"`
    BestRoute []string `json:"best_route_location_ids"`
    TotalCost int `json:"total_uber_costs"`
    TotalDuration int `json:"total_uber_duration"`
    TotalDistance float64 `json:"total_distance"`
    Eta int `json:"uber_wait_time_eta"`
}

type etaPost struct{
    Eta int `json:"eta"`
    RequestID string `json:"request_id"`
}

type TripDetails struct{
    Id  string  `json:"id"`
    Status string `json:"status"`
    StartingLocation string `json:"starting_from_location_id"`
    BestRoute []string `json:"best_route_location_ids"`
    TotalCost int `json:"total_uber_costs"`
    TotalDuration int `json:"total_uber_duration"`
    TotalDistance float64 `json:"total_distance"`
}


type TripDetailsDB struct{
    Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
    Status string `json:"status"`
    StartingLocation string `json:"starting_from_location_id"`
    BestRoute []string `json:"best_route_location_ids"`
    TotalCost int `json:"total_uber_costs"`
    TotalDuration int `json:"total_uber_duration"`
    TotalDistance float64 `json:"total_distance"`
    Index int `json:"index"`
}




type Price struct {
    ProductId       string  `json:"product_id"`
    CurrencyCode    string  `json:"currency_code"`
    DisplayName     string  `json:"display_name"`
    Estimate        string  `json:"estimate"`
    LowEstimate     int     `json:"low_estimate"`
    HighEstimate    int     `json:"high_estimate"`
    SurgeMultiplier float64 `json:"surge_multiplier"`
    Duration        int     `json:"duration"`
    Distance        float64 `json:"distance"`
}


type PlanTripInp struct {
    StartingLocation    string `json:"starting_from_location_id"`
    LocationsIds   []string `json:"location_ids"` 
}


func postPlanTrip(rw http.ResponseWriter, req *http.Request, p httprouter.Params){

    input:=PlanTripInp{}
    json.NewDecoder(req.Body).Decode(&input)

    fmt.Println(len(input.LocationsIds))
    fmt.Println(input.StartingLocation)

                //send the data to function to retrieve from DB
               // the detination too in the array in the end
    input.LocationsIds=append((input.LocationsIds),input.StartingLocation)
    m,errfromlatlongDb:=getLatLongFromDB(input.LocationsIds)

    if errfromlatlongDb!=nil{
     errorCheck(errfromlatlongDb.Error(),rw)
     return
 }

 postTripOutput,bestRouteError:=obtainRoute(m,input.LocationsIds)
 if bestRouteError!=nil{
    fmt.Print(bestRouteError.Error())
    errorCheck(bestRouteError.Error(),rw)
    return
}

respJson, errMarshall := json.Marshal(postTripOutput)
if errMarshall!=nil{
    fmt.Print("Error occcured in marshalling")
    errorCheck("Oops,something went wrong! Try after a while.",rw)
    return
}
            //sending it in the response
rw.Header().Set("Content-Type","application/json")
rw.WriteHeader(http.StatusCreated)
fmt.Fprintf(rw, "%s", respJson)    
}



         //(map[string]Coordinates,error)
func getLatLongFromDB(locationIds []string)(map[string]Coordinates,error){
    var m map[string]Coordinates
    m = make(map[string]Coordinates)

            //REPEAT CODE
              //MongoLab connection
    session,c, err := connectToDB(dbURL,locationDBName,locationCollectionName)
    if err!=nil{
        fmt.Println(err.Error())
        err_db:=errors.New("Database connection error.")
        return nil,err_db
    }
    defer session.Close()

   


    for  i:=0;i<len(locationIds);i++{
        result := LocationDBResponse{}
         //Check if given id is a valid hexademical string
                err_Hex:=checkHexString(locationIds[i])
                if err_Hex!=nil{
                     return nil,errors.New("The given location is not in the correct format")
                }
        err2 := c.Find(bson.M{"_id": bson.ObjectIdHex(locationIds[i])}).One(&result)
            //fmt.Println(result)
        if err2 != nil {
        //log.Fatal(err2)
            errMsg:=err2.Error()
            fmt.Println("inside get- error")
            if err2.Error()==ErrNotFound{
                errMsg="The given location id is incorrect. Please verify."
            }    
            return nil,errors.New(errMsg)
        }

        latLongObj:=result.Coordinate

        fmt.Println("returned lat:", result.Coordinate.Lat)
        fmt.Println("returned long:", result.Coordinate.Lng)

        m[locationIds[i]]=latLongObj
    }

            //print the value of map length
  /*  for key, value := range m {
        fmt.Println("Key:", key, "Value:", value.Lat)
    }
*/
    return m,nil

}


func getTrip(rw http.ResponseWriter, req *http.Request, p httprouter.Params){

    idString:=p.ByName("trip_id")

    //Check if given id is a valid hexademical string
    err_Hex:=checkHexString(idString)
    if err_Hex!=nil{
            //fmt.Println("The given location is not in the correct format")
       errorCheck("The given location is not in the correct format",rw)
       return
   }


        //Obtaining the json from DB

        //MongoLab connection
   session,c, err := connectToDB(dbURL,locationDBName,tripCollectionName)
   if err!=nil{
    errorCheck("Database connection error.",rw)
    return}
    defer session.Close()


    result := TripDetailsDB{}
    err2 := c.Find(bson.M{"_id": bson.ObjectIdHex(idString)}).One(&result)
    fmt.Println(result)
    if err2 != nil {
        //log.Fatal(err2)
        errMsg:=err2.Error()
        fmt.Println("inside get- error")
        if err2.Error()==ErrNotFound{
            errMsg="The given location id is incorrect. Please verify."
        }    
        errorCheck(errMsg,rw)
        return
    }

    fmt.Println("Status:", result.Status)
    fmt.Println("StartingLocation",result.StartingLocation)
    fmt.Println("BestRoute",result.BestRoute)
    fmt.Println("id",result.Id)

    //Marshling values into json
     //creating the json response
    respStruct:=TripDetails{idString,result.Status,result.StartingLocation,result.BestRoute,result.TotalCost,result.TotalDuration,result.TotalDistance}
    respJson, err4 := json.Marshal(respStruct)
    if err4!=nil{
        fmt.Print("Error occcured in marshalling")     
    }

    //sending it in the response
    rw.Header().Set("Content-Type","application/json")
    rw.WriteHeader(http.StatusOK)
    fmt.Fprintf(rw, "%s", respJson)
}


func requestTrip(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
    idString:=p.ByName("trip_id")
    fmt.Println(idString)

    //get the request id data 

    err_Hex:=checkHexString(idString)
    if err_Hex!=nil{
            //fmt.Println("The given location is not in the correct format")
       errorCheck("The given location is not in the correct format",rw)
       return
   }


        //Obtaining the json from DB
        //MongoLab connection
   session,c, err := connectToDB(dbURL,locationDBName,tripCollectionName)
   if err!=nil{
    errorCheck("Database connection error.",rw)
    return}
    defer session.Close()

    tripID:=bson.ObjectIdHex(idString)

    result := TripDetailsDB{}
    err2 := c.Find(bson.M{"_id": bson.ObjectIdHex(idString)}).One(&result)
    //fmt.Println(result)
    if err2 != nil {
        //log.Fatal(err2)
        errMsg:=err2.Error()
        fmt.Println("inside get- error")
        if err2.Error()==ErrNotFound{
            errMsg="The given location id is incorrect. Please verify."
        }    
        errorCheck(errMsg,rw)
        return
    }

    fmt.Println("Status:", result.Status)
    fmt.Println("StartingLocation",result.StartingLocation)
    fmt.Println("BestRoute",result.BestRoute)
    fmt.Println("id",result.Id)


    //CHECK if the trip is completed
    if(result.Status=="completed"){
        // var dummyArray []string
       completedStruct:=PutTrip{idString,result.Status,result.StartingLocation,"",result.BestRoute,result.TotalCost,result.TotalDuration,result.TotalDistance,0}

    //Send the response json:
       putResponseJson, err4 := json.Marshal(completedStruct)
       if err4!=nil{
        fmt.Print("Error occcured in marshalling")     
    }

    //sending it in the response
    rw.Header().Set("Content-Type","application/json")
    rw.WriteHeader(http.StatusCreated)
    fmt.Fprintf(rw, "%s", putResponseJson)


}else{


    //Getting the best route array and source place

       //Getting the source index 
   startingLocID,currDestID,statustoSetDB,err_cal:=SourceDestinationDetermination(result)
   if err_cal!=nil{
       fmt.Println("Error in calculation of source and dest")
   }

    //get the Coordinates of the 2 of them
   session1,c1, err1 := connectToDB(dbURL,locationDBName,locationCollectionName)
   if err1!=nil{
    errorCheck("Database connection error.",rw)
    return}
    defer session1.Close()

    //Co-ordinate 1
    resultCord1 := LocationDBResponse{}
    err3 := c1.Find(bson.M{"_id": bson.ObjectIdHex(startingLocID)}).One(&resultCord1)
    fmt.Println(resultCord1)
    if err3 != nil {
        //log.Fatal(err2)
        errMsg1:=err3.Error()
        fmt.Println("inside get- error")
        if err3.Error()==ErrNotFound{
            errMsg1="The given location id is incorrect. Please verify."
        }    
        errorCheck(errMsg1,rw)
        return
    }
    startingCoord:=resultCord1.Coordinate
    

    //destn coord
    err_dest := c1.Find(bson.M{"_id": bson.ObjectIdHex(currDestID)}).One(&resultCord1)
    fmt.Println(resultCord1)
    if err_dest != nil {
        //log.Fatal(err2)
        errMsg1:=err_dest.Error()
        fmt.Println("inside get- error")
        if err_dest.Error()==ErrNotFound{
            errMsg1="The given location id is incorrect. Please verify."
        }    
        errorCheck(errMsg1,rw)
        return
    }
    destCoord:=resultCord1.Coordinate

    //Printing the coordinates:
    fmt.Println("Coordinates-src: ",startingCoord.Lat,startingCoord.Lng,"Coordinates-dest: ",destCoord.Lat,destCoord.Lng)

//Uber call with current source and destination
    strStartLat:=strconv.FormatFloat(startingCoord.Lat,'f',7,64)
    strStartLon:=strconv.FormatFloat(startingCoord.Lng,'f',7,64)
    strEndLat:=strconv.FormatFloat(destCoord.Lat,'f',7,64)
    strEndLon:=strconv.FormatFloat(destCoord.Lng,'f',7,64)


    //url:="https://sandbox-api.uber.com/v1/requests?start_latitude=37.625732&start_longitude=-122.377807&end_latitude=37.785114&end_longitude=-122.406677&product_id=a1111c8c-c720-46c3-8534-2fcdd730040d"
    url:="https://sandbox-api.uber.com/v1/requests?start_latitude="+strStartLat+"&start_longitude="+strStartLon+"&end_latitude="+strEndLat+"&end_longitude="+strEndLon+"&product_id=2832a1f5-cfc0-48bb-ab76-7ea7a62060e7"
    fmt.Println(url)



    var jsonStr = []byte(`{"start_longitude":"`+strStartLon+`","product_id":"2832a1f5-cfc0-48bb-ab76-7ea7a62060e7","start_latitude":"`+strStartLat+`"}`)
    //fmt.Println(jsonStr)
    req1, errReqC := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    if errReqC!=nil{
        errMsg:="Uber Server/Api is not responding. Please check after a while"
         errorCheck(errMsg,rw)
        return
    }
    //TODO: Check for errReqC

    req1.Header.Set("Content-Type", "application/json")
    req1.Header.Set("Authorization", authToken)

    client := &http.Client{}
    resp, erruber := client.Do(req1)
    if erruber != nil {
        errMsg:="Uber Server/Api is not responding. Please check after a while"
        errorCheck(errMsg,rw)
        return
    }
    defer resp.Body.Close()

    //Update the DB entry with the status and id
    c.UpdateId(tripID,bson.M{"$set": bson.M{"status": statustoSetDB,"index":result.Index+1}})



   // fmt.Println("response Status:", resp.Status)
   // fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))

    //Form the PUT return json
    var etaStruct etaPost
    err_unmarshal:=json.Unmarshal(body,&etaStruct)
    if err_unmarshal!=nil{
        fmt.Println("Unmarshalling error at eta: ", err_unmarshal)
        errorCheck("Oops,something went wrong! Try after a while.",rw)
        return
   }

    //got the eta
   fmt.Println("Obtained Eta: ", etaStruct.Eta)
    //eta:=etaStruct.Eta

    //Get eta from uber PUT 
   eta,err_eta:=getETAFromPut(etaStruct.RequestID)
   if err_eta!=nil{
        errorCheck(err_eta.Error(),rw)
        return   
}

    //modified array
    //newBestRoute:=result.BestRoute[result.Index:]

putResponseStruct:=PutTrip{idString,"requesting",result.StartingLocation,currDestID,result.BestRoute,result.TotalCost,result.TotalDuration,result.TotalDistance,eta}

    //Send the response json:
putResponseJson, err4 := json.Marshal(putResponseStruct)
if err4!=nil{
    fmt.Print("Error occcured in marshalling")  
     errorCheck("Oops,something went wrong! Try after a while.",rw)
        return   
}

    //sending it in the response
rw.Header().Set("Content-Type","application/json")
rw.WriteHeader(http.StatusCreated)
fmt.Fprintf(rw, "%s", putResponseJson)
}


}

func getETAFromPut(reqID string) (int,error){
    //PUT to UBER api
    url:="https://sandbox-api.uber.com/v1/sandbox/requests/"+reqID
    fmt.Println(url)
    var jsonStr = []byte(`{"status: "accepted"}`)
    fmt.Println(jsonStr)
    req1, errEta := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
    //TODO: Check for errReqC
    if errEta!=nil{
        fmt.Println(errEta.Error())
        return 0,errors.New("Http Request creation issue")
    }

    req1.Header.Set("Content-Type", "application/json")
    req1.Header.Set("Authorization", authToken)

    client := &http.Client{}
    resp, erruber := client.Do(req1)
    if erruber != nil {
       fmt.Println(erruber.Error())
       return 0,errors.New("Uber api request issue!")
   }
   defer resp.Body.Close()



    //GET the eta using uber api
   urlGetLink:="https://sandbox-api.uber.com/v1/requests/"+reqID
   fmt.Println(urlGetLink)
   reqGet, errGetReq := http.NewRequest("GET", urlGetLink, nil)
   if errGetReq!=nil{
    fmt.Println(errGetReq.Error())
    return 0,errors.New("Http Request creation issue")
}
reqGet.Header.Set("Authorization", authToken)
clientGet := &http.Client{}
respGet, errGet := clientGet.Do(reqGet)
if errGet != nil {
    fmt.Println(errGet.Error())
    err_reqUI:=errors.New("Oops,something went wrong! Try after a while.")
    return 0,err_reqUI
}
defer respGet.Body.Close()

body, err1 := ioutil.ReadAll(respGet.Body)
if err1 != nil {
    return 0,errors.New("Oops,something went wrong! Try after a while.")

}
var etaStruct etaPost 
                //Unmarshall the response into a json
err2:=json.Unmarshal(body,&etaStruct)
if err2 != nil {
    return 0,errors.New("Oops,something went wrong! Try after a while.")
}

fmt.Println("Obtained put req id and eta",etaStruct.Eta, etaStruct.RequestID)
return etaStruct.Eta,nil




}


func SourceDestinationDetermination(dbResult TripDetailsDB)(string,string,string,error){
    index:=dbResult.Index
    routeArray:=dbResult.BestRoute

    //First element
    if(index==0){
        return dbResult.StartingLocation,routeArray[index],"requesting",nil
    }        //If last element
    if(index==len(routeArray)){
        return routeArray[index-1],dbResult.StartingLocation,"completed",nil
    } else{
        //Otherwise
       return routeArray[index-1],routeArray[index],"requesting",nil 
   }

    //use try and catch Block


}


         //


func getUberdata(StartLatitude float64, StartLongitude float64, EndLatitude float64, EndLongitude float64) (Price,error){

            //link:="https://api.uber.com/v1/estimates/price?start_latitude="+strconv.FormatFloat(StartLongitude,'f',2,64)+"&start_longitude="+strconv.FormatFloat(StartLongitude,'f',2,31)+"&end_latitude="+strconv.FormatFloat(EndLatitude,'f',2,31)+"&end_longitude="+strconv.FormatFloat(EndLongitude,'f',2,31)+"&server_token=fN-83tP1-jl7M9fxfvsJYSVQbXZhEGqUsIxDdIRK"
    strStartLat:=strconv.FormatFloat(StartLatitude,'f',7,64)
    strStartLon:=strconv.FormatFloat(StartLongitude,'f',7,64)
    strEndLat:=strconv.FormatFloat(EndLatitude,'f',7,64)
    strEndLon:=strconv.FormatFloat(EndLongitude,'f',7,64)

    //link:="https://api.uber.com/v1/estimates/price?start_latitude="+strStartLat+"&start_longitude="+strStartLon+"&end_latitude="+strEndLat+"&end_longitude="+strEndLon+"&server_token=fN-83tP1-jl7M9fxfvsJYSVQbXZhEGqUsIxDdIRK"
    link:="https://sandbox-api.uber.com/v1/estimates/price?start_latitude="+strStartLat+"&start_longitude="+strStartLon+"&end_latitude="+strEndLat+"&end_longitude="+strEndLon+"&server_token=fN-83tP1-jl7M9fxfvsJYSVQbXZhEGqUsIxDdIRK"
                //fmt.Println(link)
    dummy:=Price{}

    resp, err := http.Get(link);
    if err != nil {
        fmt.Println(err.Error())
        err_google:=errors.New("Uber map api connection could not be established!")
        return dummy,err_google
    }


    defer resp.Body.Close()
    body, err1 := ioutil.ReadAll(resp.Body)
    if err1 != nil {
        //fmt.Println("Error: ", err1.Error())
        return dummy,err1

    }
    var otainedPrices PriceEstimatesUber
                //Unmarshall the response into a json
    err2:=json.Unmarshal(body,&otainedPrices)
    if err2 != nil {
        return dummy,err2
    }

                    //TODO: incorrect lat long, server down etc

    //if error, we dont get Prices array is 0
    if((len(otainedPrices.Prices))==0){
        return dummy,errors.New("Invalid data of the location pair: could be latitude/longitude ")
    }
       
    lowestPrice:=checkLowestEstimatePrice(otainedPrices)
    fmt.Println(lowestPrice.LowEstimate)

    return lowestPrice,nil


}


func obtainRoute(totalMap map[string]Coordinates,locationIds []string)(TripDetails,error){
    m :=totalMap
    locIDs:=locationIds
    dummyTripDetails:=TripDetails{}


            //remove from location and totalMap and put in source variable
    var sourceCoor Coordinates  = m[locIDs[len(locIDs)-1]]
    originalSourceCoord := sourceCoor
    var sourceId string =locIDs[len(locIDs)-1]

    var total_uber_costs int
    var total_uber_duration int
    var total_distance float64


    fmt.Println("Obtained source"+sourceId)
            //deleting from map and array list
    delete(m,sourceId)
    lastindex := len(locIDs)-1
    fmt.Println("length befor: "+strconv.Itoa(len(locIDs)))
    locIDs=append(locIDs[:lastindex])
    fmt.Println("length after: "+strconv.Itoa(len(locIDs)))


    fmt.Println("*************************")
             //print the value of map length
    fmt.Println(sourceCoor.Lat)

    for key, value := range m {
        fmt.Println("Key:", key, "Value:", value.Lat)
    }


            //map of destinatins and dist,duration between Source and destinations
    var totalPairMap map [string][]int
    totalPairMap = make(map [string][]int)

            //datastructure for holding the order of destinations
    var locationShortestDist []string




            //Take the pairs 
            //First round 


    noOfDestinations:=len(locIDs)
    intialDestinations:=len(locIDs)



    for j:=0;j<intialDestinations;j++{

        fmt.Println("###################################################")
        fmt.Println("start of iteration: ",j)
        fmt.Println("noOfDestinations: ",noOfDestinations)
        fmt.Println("source latitude: ",sourceCoor.Lat)

        var tempLowestCost int
        var lowID string  
        var lowIndex int
        var tempLowestDuration int



        var tempPriceMaps map[string]Price
        tempPriceMaps = make(map[string]Price) 



        for i:=0;i<noOfDestinations;i++{

            var individualPrice Price
            var ubrErr error

            next:=m[locIDs[i]]
            individualPrice,ubrErr=getUberdata(sourceCoor.Lat,sourceCoor.Lng,next.Lat,next.Lng)
            if(ubrErr!=nil){
                return dummyTripDetails,ubrErr
            }

            tempPriceMaps[locIDs[i]]=individualPrice

                            //Print lowestimate, duration, cost, 

                fmt.Println("Key:",locIDs[i],"Lat:",next.Lat,"LE",individualPrice.LowEstimate,"Dur",individualPrice.Duration,"Dis",individualPrice.Distance)

                 //if last item, add the element and exit
                if(i==0){
                 lowID=locIDs[i]
                 tempLowestCost=individualPrice.LowEstimate
                 tempLowestDuration=individualPrice.Duration}




                //fmt.Println(individualPrice.LowEstimate)
                //inserting the values into a map
                 totalPairMap[locIDs[i]]=[]int{individualPrice.LowEstimate,individualPrice.Duration}
                //Check for LowEstimate
                 if(individualPrice.LowEstimate<tempLowestCost){
                    //TODO check for equality
                    tempLowestCost=individualPrice.LowEstimate
                    tempLowestDuration=individualPrice.Duration
                    lowID=locIDs[i]
                    lowIndex=i
                }else if(individualPrice.LowEstimate==tempLowestCost){
                    //now compare with duration
                    if(individualPrice.Duration<tempLowestDuration){
                    //tempLowestCost=individualPrice.LowEstimate
                        tempLowestDuration=individualPrice.Duration
                        lowID=locIDs[i]
                        lowIndex=i
                    }
                }




                       // fmt.Println("map of distance and cost")

                      /*  for key, value := range totalPairMap {
                fmt.Println("Key:", key, "Value:", value[0],value[1])*/



            }

            //Adding to the  final list
            fmt.Println("The id with lowest cost:"+lowID)
            locationShortestDist = append(locationShortestDist,lowID)

                //Cumulating everything
            total_uber_costs += tempPriceMaps[lowID].LowEstimate
            total_uber_duration += tempPriceMaps[lowID].Duration
            total_distance += tempPriceMaps[lowID].Distance

            fmt.Println("Total distance each time: ", total_distance)

            //deleting from the map, the lowest entry
            delete(totalPairMap,lowID)
            //making the destination one lesser
            noOfDestinations=noOfDestinations-1
            //Removing from the array index

                //check if it is the last element in the array:
            fmt.Println("lowIndex ",lowIndex)
            fmt.Println("len(locIDs)",len(locIDs))


            if(lowIndex==(len(locIDs)-1)){
                locIDs = locIDs[:len(locIDs)-1]
                }else{locIDs = append(locIDs[:lowIndex], locIDs[lowIndex+1:]...)}

                //Changing the source:
                sourceCoor  = m[lowID]




            //Re printing the map
                /*fmt.Println("Re printing the map of distance and cost")
                for key, value := range totalPairMap {
                    fmt.Println("Key:", key, "Value:", value[0],value[1])

                    }*/

                }

                //Adding the last destination to source
                lastDEs:=totalMap[locationShortestDist[len(locationShortestDist)-1]]
                fmt.Println("Last dest: ", lastDEs.Lat,lastDEs.Lng)

                //now uber call with last dest and source
                //originalSourceCoord
                individualPrice1,ubrErr1:=getUberdata(lastDEs.Lat,lastDEs.Lng,originalSourceCoord.Lat,originalSourceCoord.Lng)
                if ubrErr1!=nil{
                    //fmt.Println("Uber data retrival issue",ubrErr1.Error())
                    return dummyTripDetails,ubrErr1
                }

                //Add to the total cost, distance and duration
                total_uber_costs+=individualPrice1.LowEstimate
                total_uber_duration+=individualPrice1.Duration
                total_distance+=individualPrice1.Distance



                //print the arraylist of the locations added:
              /*  for value1:= range locationShortestDist{
                fmt.Println("value",value1)}*/

                fmt.Println("length of final array",len(locationShortestDist)," first element",locationShortestDist[0])

                for k:=0;k<len(locationShortestDist);k++{
                    fmt.Println("val",locationShortestDist[k])
                }

                fmt.Println(total_uber_costs,total_uber_duration,total_distance)


                //Adding to the collection
                session,c, err := connectToDB(dbURL,locationDBName,tripCollectionName)
                if err!=nil{
                    fmt.Println("Database connection error.")
                    //errorCheck("Database connection error.",rw)
                    return dummyTripDetails,errors.New("Database error!")}
                    defer session.Close()
                    


                    i := bson.NewObjectId()
                    //fmt.Println(i)
                    //fmt.Println("String version")
                    idString:=i.String()
                    //fmt.Println(idString)

                             //Extracting the ID
                    r, err := regexp.Compile(`"[a-z0-9]+"`)

                if err != nil {
            
                     return dummyTripDetails, errors.New("Oops,something went wrong! Try after a while.")

                 }

                    split1:=r.FindString(idString)
                    ID:=strings.Trim(split1,"\"")



                    //first time entry into the DB
                    d:= TripDetailsDB{i,"planning",sourceId,locationShortestDist,total_uber_costs,total_uber_duration,total_distance,0}
                    err_ins:=c.Insert(d)
                     if err_ins != nil {
                    fmt.Println("Database error")
                     return dummyTripDetails, errors.New("Database Error. Try after a while.")

                 }
                    //creation of the json response
                    respStruct:=TripDetails{ID,"planning",sourceId,locationShortestDist,total_uber_costs,total_uber_duration,total_distance}
                 return respStruct, nil
             }        


            //Function to get the lowest of the several estimates obtained from uber for a pair of locations
             func checkLowestEstimatePrice(totalPE PriceEstimatesUber) Price {
            //var lowestPrice Price
                priceArray:=totalPE.Prices
                temp := priceArray[0].LowEstimate
                lowIndex :=0


                for i:=1;i<len(priceArray);i++{
                //Checking for null values
                    if(priceArray[i].LowEstimate!=0 && priceArray[i].DisplayName!="Health"){
                        if(priceArray[i].LowEstimate<temp){
                            temp=priceArray[i].LowEstimate
                            lowIndex=i
                        }
                    }
                }
                //return lowest
                fmt.Println(priceArray[lowIndex].DisplayName)
                return priceArray[lowIndex]

            }





