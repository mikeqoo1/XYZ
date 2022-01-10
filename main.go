package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type 原始資料 struct {
	ttime string
	data  string
	流口水號  string
}

type 原始資料List []原始資料

func (c 原始資料List) Len() int      { return len(c) }
func (c 原始資料List) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c 原始資料List) Less(i, j int) bool {
	if c[i].ttime == c[j].ttime {
		return c[i].流口水號 < c[j].流口水號
	} else {
		return c[i].ttime < c[j].ttime
	}
} // 小到大排序

func 合併檔案(newfile *os.File, oldfile *os.File) {
	n, err := io.Copy(newfile, oldfile)
	if err != nil {
		fmt.Println("failed to append signed file to output=", err)
	}
	fmt.Printf("wrote %d bytes \n", n)
}

func main() {
	/*排序:https://segmentfault.com/a/1190000008062661*/

	// open the file
	out, err := os.OpenFile("File/AllData.txt", os.O_CREATE|os.O_WRONLY, 0755)
	//handle errors while opening
	if err != nil {
		fmt.Println("Error when opening file:", err)
	}

	EMSTorder1, _ := os.Open("File/檔案1") //emst order
	合併檔案(out, EMSTorder1)
	EMSTorder2, _ := os.Open("File/檔案2") //emst order
	合併檔案(out, EMSTorder2)
	EMST基本價, _ := os.Open("File/檔案3") //c67
	合併檔案(out, EMST基本價)
	ben, _ := os.Open("File/檔案4") // ben order
	合併檔案(out, ben)
	pvc, _ := os.Open("File/檔案5") //pvc log
	合併檔案(out, pvc)
	EMST成交, _ := os.Open("File/檔案6") //emst成交回報
	合併檔案(out, EMST成交)
	EMST委託, _ := os.Open("File/檔案7") //emst委託回報
	合併檔案(out, EMST委託)
	Want, _ := os.Open("File/wamt.log")
	合併檔案(out, Want)
	out.Close()

	//這邊要重新開檔 因為上面我只設定寫的權限 func NewScanner 要讀
	file, _ := os.Open("File/AllData.txt")
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	All := []原始資料{} //排序用
	// read line by line
	for fileScanner.Scan() {
		var temp 原始資料
		if strings.Contains(fileScanner.Text(), "stock_id") {
			temp.ttime = "000000000"
			temp.data = ""
		} else if len(fileScanner.Text()) == 120 { //紀錄電檔案1文長度120(emst)下單
			動作別 := fileScanner.Text()[10:11]
			if 動作別 == "I" {
				//temp.ttime = fileScanner.Text()[26:35]
				temp.ttime = fileScanner.Text()[26:28] + ":" + fileScanner.Text()[28:30] + ":" + fileScanner.Text()[30:32] + "." + fileScanner.Text()[32:35]
				temp.data = "8988 |" + temp.ttime + "|" + fileScanner.Text()
			} else if 動作別 == "C" || 動作別 == "P" {
				//fmt.Println(fileScanner.Text()[25:34])
				temp.ttime = fileScanner.Text()[25:27] + ":" + fileScanner.Text()[27:29] + ":" + fileScanner.Text()[29:31] + "." + fileScanner.Text()[31:34]
				temp.data = "8988 |" + temp.ttime + "|" + fileScanner.Text()
			} else if 動作別 == "D" {
				//temp.ttime = fileScanner.Text()[25:34]
				temp.ttime = fileScanner.Text()[25:27] + ":" + fileScanner.Text()[27:29] + ":" + fileScanner.Text()[29:31] + "." + fileScanner.Text()[31:34]
				temp.data = "8988 |" + temp.ttime + "|" + fileScanner.Text()
			}
		} else if len(fileScanner.Text()) == 108 { //emst成交
			HH := fileScanner.Text()[75:77]
			MM := fileScanner.Text()[77:79]
			SS := fileScanner.Text()[79:81]
			sss := fileScanner.Text()[81:84]
			i, _ := strconv.Atoi(SS)
			i = i + 6
			if i >= 60 {
				l, _ := strconv.Atoi(MM)
				l = l + 1
				MM = strconv.Itoa(l)
				if len(MM) < 2 {
					MM = "0" + MM
				}
				SS = strconv.Itoa(i - 60)
				if len(SS) < 2 {
					SS = "0" + SS
				}
			} else {
				SS = strconv.Itoa(i)
				if len(SS) < 2 {
					SS = "0" + SS
				}
				sss = "000"
			}
			temp.ttime = HH + ":" + MM + ":" + SS + "." + sss
			temp.data = "55688|" + temp.ttime + "|" + fileScanner.Text()
		} else if len(fileScanner.Text()) == 162 { //emst委託
			HH := fileScanner.Text()[71:73]
			MM := fileScanner.Text()[73:75]
			SS := fileScanner.Text()[75:77]
			sss := fileScanner.Text()[77:80]
			i, _ := strconv.Atoi(SS)
			i = i + 5
			if i >= 60 {
				l, _ := strconv.Atoi(MM)
				l = l + 1
				MM = strconv.Itoa(l)
				if len(MM) < 2 {
					MM = "0" + MM
				}
				SS = strconv.Itoa(i - 60)
				if len(SS) < 2 {
					SS = "0" + SS
				}
			} else {
				SS = strconv.Itoa(i)
				if len(SS) < 2 {
					SS = "0" + SS
				}
				sss = "000"
			}
			temp.ttime = HH + ":" + MM + ":" + SS + "." + sss
			temp.data = "55688|" + temp.ttime + "|" + fileScanner.Text()
		} else if strings.Contains(fileScanner.Text(), "order") { //Ben
			qqq := strings.Split(fileScanner.Text(), ",")
			temp.ttime = qqq[0]
			if qqq[2] == "1t" {
				temp.data = "3333 |" + temp.ttime + "|" + "\x01" + qqq[3]
			} else if qqq[2] == "1s" {
				temp.data = "13334|" + temp.ttime + "|" + "\x01" + qqq[3]
			} else {
				temp.data = "13335|" + temp.ttime + "|" + "\x01" + qqq[3]
			}
		} else if strings.Contains(fileScanner.Text(), "FIX.4.4") || strings.Contains(fileScanner.Text(), "TMP_") {
			ooo := strings.Split(fileScanner.Text(), ",")
			temp.ttime = ooo[0]
			if ooo[1] == "FIX_O" || ooo[1] == "FIX_T" {
				長度 := len(ooo[2])
				www := strings.Split(ooo[2], "\x01")
				for i := 0; i < len(www); i++ {
					if strings.Contains(www[i], "34=") {
						temp.流口水號 = www[i][3:]
					}
				}
				mm := temp.ttime[6:8]
				sss := temp.ttime[9:]
				i, _ := strconv.Atoi(mm)
				i = i + 1
				if i >= 60 {
					mm = temp.ttime[6:8]
					sss = "999"
				} else {
					mm = strconv.Itoa(i)
					if len(mm) < 2 {
						mm = "0" + mm
					}
					sss = temp.ttime[9:]
				}
				temp.ttime = temp.ttime[0:5] + ":" + mm + "." + sss
				temp.data = "55689|" + temp.ttime + "|" + ooo[2][:長度-8]
				temp.data = strings.ReplaceAll(temp.data, "body為:", "")
			} else if ooo[1] == "TMP_O" || ooo[1] == "TMP_T" {
				mm := temp.ttime[6:8]
				sss := temp.ttime[9:]
				i, _ := strconv.Atoi(mm)
				i = i + 2
				if i >= 60 {
					mm = temp.ttime[6:8]
					sss = "999"
				} else {
					mm = strconv.Itoa(i)
					if len(mm) < 2 {
						mm = "0" + mm
					}
					sss = temp.ttime[9:]
				}
				temp.ttime = temp.ttime[0:5] + ":" + mm + "." + sss
				temp.data = "55690|" + temp.ttime + "|" + ooo[2]
			}
		} else if len(fileScanner.Text()) > 100 && fileScanner.Text()[84:85] == "E" {
			temp.data = "風控Error|" + fileScanner.Text()
			fmt.Println(temp.data)
		} else if strings.Contains(fileScanner.Text(), "WAMT") {
			ooo := strings.Split(fileScanner.Text(), ",")
			temp.ttime = ooo[0][9:] + ".000"
			a := strings.Index(fileScanner.Text(), "'")
			//fmt.Println(strings.Index(fileScanner.Text(), "'"))
			b := strings.LastIndex(fileScanner.Text(), "'")
			//fmt.Println(strings.LastIndex(fileScanner.Text(), "'"))
			ddd := fileScanner.Text()[a+1 : b]
			//fmt.Println(fileScanner.Text()[a:b])
			temp.data = "3306 |" + temp.ttime + "|" + ddd
			//fmt.Println(temp.ttime)
			//fmt.Println(temp.data)
		} else { //price
			temp.ttime = fileScanner.Text()[11:19]
			temp.ttime = temp.ttime + ".000"
			temp.data = "8986 |" + temp.ttime + "|" + fileScanner.Text()
		}
		All = append(All, temp)
	}

	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		fmt.Println("Error while reading file:", err)
	}

	//排序
	sort.Sort(原始資料List(All))

	//寫檔
	write, _ := os.OpenFile("End/OnePiece.txt", os.O_CREATE|os.O_WRONLY, 0755)
	defer write.Close()
	for i := 0; i < len(All); i++ {
		if All[i].data == "" || strings.Contains(All[i].data, "風控Error") {
			//不寫檔
		} else {
			write.Write([]byte(All[i].data + "\n"))
		}
	}
	fmt.Printf("END!\n")
}
