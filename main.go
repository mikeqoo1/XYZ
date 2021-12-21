package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type 原始資料 struct {
	ttime string
	data  string
}

type 原始資料List []原始資料

func (c 原始資料List) Len() int           { return len(c) }
func (c 原始資料List) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c 原始資料List) Less(i, j int) bool { return c[i].ttime < c[j].ttime } // 小到大排序

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
	out, err := os.OpenFile("AllData.txt", os.O_CREATE|os.O_WRONLY, 0755)
	//handle errors while opening
	if err != nil {
		fmt.Println("Error when opening file:", err)
	}

	EMSTorder1, _ := os.Open("檔案1") //emst order
	合併檔案(out, EMSTorder1)
	EMSTorder2, _ := os.Open("檔案2") //emst order
	合併檔案(out, EMSTorder2)
	EMST基本價, _ := os.Open("檔案3") //c67
	合併檔案(out, EMST基本價)
	ben, _ := os.Open("檔案4") // ben order
	合併檔案(out, ben)
	pvc, _ := os.Open("檔案5") //pvc log
	合併檔案(out, pvc)
	EMST成交, _ := os.Open("檔案6") //emst成交回報
	合併檔案(out, EMST成交)
	EMST委託, _ := os.Open("檔案7") //emst委託回報
	合併檔案(out, EMST委託)

	out.Close()

	//這邊要重新開檔 因為上面我只設定寫的權限 func NewScanner 要讀
	file, _ := os.Open("AllData.txt")
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
			temp.ttime = fileScanner.Text()[75:84]
			temp.data = "55688|" + temp.ttime + "|" + fileScanner.Text()
		} else if len(fileScanner.Text()) == 162 { //emst委託
			temp.ttime = fileScanner.Text()[71:80]
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
			if ooo[1] == "fix下單" {
				//temp.data = "55688|" + ooo[2]
			} else if ooo[1] == "FIX_O" || ooo[1] == "FIX_T" {
				temp.data = "55689|" + temp.ttime + "|" + ooo[2]
				temp.data = strings.ReplaceAll(temp.data, "body為:", "")
			} else if ooo[1] == "TMP_O" || ooo[1] == "TMP_T" {
				temp.data = "55690|" + temp.ttime + "|" + ooo[2]
			}
		} else if len(fileScanner.Text()) > 100 && fileScanner.Text()[84:85] == "E" {
			temp.data = "風控Error|" + fileScanner.Text()
			fmt.Println(temp.data)
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
