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

	EMSTorder1, _ := os.Open("檔案1")
	合併檔案(out, EMSTorder1)
	EMSTorder2, _ := os.Open("檔案2")
	合併檔案(out, EMSTorder2)
	EMST基本價, _ := os.Open("檔案3")
	合併檔案(out, EMST基本價)
	ben, _ := os.Open("檔案4")
	合併檔案(out, ben)
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
		} else if len(fileScanner.Text()) == 120 { //紀錄電檔案1文長度120(emst)
			動作別 := fileScanner.Text()[10:11]
			if 動作別 == "I" {
				temp.ttime = fileScanner.Text()[26:35]
				temp.data = "8988 |" + fileScanner.Text()
			} else if 動作別 == "C" || 動作別 == "P" {
				temp.ttime = fileScanner.Text()[25:34]
				temp.data = "8988 |" + fileScanner.Text()
			} else if 動作別 == "D" {
				temp.ttime = fileScanner.Text()[25:34]
				temp.data = "8988 |" + fileScanner.Text()
			}
		} else if strings.Contains(fileScanner.Text(), "order") { //Ben
			temp.ttime = fileScanner.Text()[0:12]
			temp.ttime = strings.ReplaceAll(temp.ttime, ":", "")
			temp.ttime = strings.ReplaceAll(temp.ttime, ".", "")
			qqq := strings.Split(fileScanner.Text(), ",")
			if qqq[2] == "1t" {
				temp.data = "3333 |" + "\x01" + qqq[3]
			} else if qqq[2] == "1s" {
				temp.data = "13334|" + "\x01" + qqq[3]
			} else {
				temp.data = "13335|" + "\x01" + qqq[3]
			}
		} else { //price
			temp.ttime = fileScanner.Text()[11:19]
			temp.ttime = strings.ReplaceAll(temp.ttime, ":", "") + "000"
			temp.data = "8986 |" + fileScanner.Text()
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
		if All[i].data == "" {
			//不寫檔
		} else {
			write.Write([]byte(All[i].data + "\n"))
		}
	}
	fmt.Printf("END!\n")
}
