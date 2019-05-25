package main

import (
  "fmt"
  "io/ioutil"
  "strconv"
  "gopkg.in/benweidig/cli-table.v2"
)

//Sonde Object to parse incoming data
type Sonde struct {
  ID        int64
  ChanNum   int64
  FrameNum  int64
  Temp      int64
  Pressure  int64
  Humidity  int64
}

//Start Symbol used to synchronise the stream.
var start  = "0000011100010010"
//Buffer is string used as a FIFO Queue
var buffer = "0000000000000000"


//Convert incoming bit as char to a string value
func convert2bin(l []byte) []string {
  v := []string{}
  for _,value := range l {
    if value == 48 {
      v = append(v,"0")
    } else {
      v = append(v,"1")
    }
  }
  return v
}

//Parse data to a Sonde Object
func Parse(bl []string) Sonde {
  sondeData := Sonde{}
  sync := 0
  for _,value := range bl {
    //Shift data into the FIFO Queue
    buffer = fmt.Sprintf(buffer[1:]+value)
    //Check if the start symbol is finded
    if buffer == start {
      //Reset sync pointer
      sync = -1
    }
    //Parse frame with scheme from this datasheet
    //http://www.meteo-tech.co.il/ImagesDownloadFiles/w9000.pdf
    if sync++; sync%16 == 0 {
      switch sync/16 {
        case 1:
          sondeData.ID, _ = strconv.ParseInt(buffer, 2, 64)
        case 2:
          sondeData.ChanNum, _ = strconv.ParseInt(buffer, 2, 64)
        case 3:
          sondeData.FrameNum, _ = strconv.ParseInt(buffer, 2, 64)
        case 5:
          sondeData.Temp, _ = strconv.ParseInt(buffer, 2, 64)
        case 6:
          sondeData.Pressure, _ = strconv.ParseInt(buffer, 2, 64)
        case 7:
          sondeData.Humidity, _ = strconv.ParseInt(buffer, 2, 64)
        }
    }
  }
  return sondeData
}


func DisplayData(s Sonde){
  //Create date table
  table := clitable.New()
  table.ColSeparator = ":"
  table.AddRow("ID ",s.ID)
  table.AddRow("Number of channels ",s.ChanNum)
  table.AddRow("Frame count number ",s.FrameNum)
  table.AddRow("Temperature ",s.Temp)
  table.AddRow("Pressure ",s.Pressure)
  table.AddRow("Humidity ",s.Humidity)
  //Diplay
  fmt.Println("")
  fmt.Println("Data From an Unknown Meteorological Station :")
	fmt.Println(table.String())

}

func main(){
  f, _ := ioutil.ReadFile("demodded.txt")
  incomingData := Parse(convert2bin(f))
  DisplayData(incomingData)
}
