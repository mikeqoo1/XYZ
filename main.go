package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
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
	盤中, _ := os.Open("File/盤中異動")
	合併檔案(out, 盤中)

	// # rename $1 $2 $3
	// # $1: 要被取代的關鍵字
	// # $2: 新的關鍵字
	// # $3: 檔名符合這個規則的才取代
	//去236 /home/Projects/file 執行 rename 01110105 01110118 *.TXT.*

	// 	雖然 sed 的選項和指令種類繁多，但實務上經常用到的大致有以下幾種：
	// -n：沉默模式。
	// -e：直接在命令模式編輯。(可不加，請見詳解)
	// -f：程式手稿不直接在命列中打上，而是從指定的檔案中載入。
	// -i：修改檔案。
	//去236 /home/Projects/file 執行 sed -i 's/01110105/01110118/g' *

	out.Close()

	//這邊要重新開檔 因為上面我只設定寫的權限 func NewScanner 要讀
	file, _ := os.Open("File/AllData.txt")
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	All := []原始資料{} //排序用
	// read line by line
	for fileScanner.Scan() {
		var temp 原始資料
		month := time.Now().Month()
		day := time.Now().Day()
		資料內容 := fileScanner.Text()
		if strings.Contains(fileScanner.Text(), "20220105") {
			std := fmt.Sprintf("2022%02d%02d", int(month), day)
			資料內容 = strings.Replace(fileScanner.Text(), "20220105", std, -1)
		}
		if strings.Contains(fileScanner.Text(), "01110105") {
			std := fmt.Sprintf("0111%02d%02d", month, day)
			資料內容 = strings.Replace(fileScanner.Text(), "01110105", std, -1)
		}
		if strings.Contains(fileScanner.Text(), "1110105") {
			std := fmt.Sprintf("111%02d%02d", month, day)
			資料內容 = strings.Replace(fileScanner.Text(), "1110105", std, -1)
		}
		if strings.Contains(資料內容, "stock_id") {
			temp.ttime = "000000000"
			temp.data = ""
		} else if len(資料內容) == 120 { //紀錄電檔案1文長度120(emst)下單
			動作別 := 資料內容[10:11]
			if 動作別 == "I" {
				//temp.ttime = 資料內容[26:35]
				temp.ttime = 資料內容[26:28] + ":" + 資料內容[28:30] + ":" + 資料內容[30:32] + "." + 資料內容[32:35]
				temp.data = "8988 |" + temp.ttime + "|" + 資料內容
			} else if 動作別 == "C" || 動作別 == "P" {
				//fmt.Println(資料內容[25:34])
				temp.ttime = 資料內容[25:27] + ":" + 資料內容[27:29] + ":" + 資料內容[29:31] + "." + 資料內容[31:34]
				temp.data = "8988 |" + temp.ttime + "|" + 資料內容

			} else if 動作別 == "D" {
				//temp.ttime = 資料內容[25:34]
				temp.ttime = 資料內容[25:27] + ":" + 資料內容[27:29] + ":" + 資料內容[29:31] + "." + 資料內容[31:34]
				temp.data = "8988 |" + temp.ttime + "|" + 資料內容
			}
		} else if len(資料內容) == 108 { //emst成交
			if 資料內容[4:5] == "V" {
				temp.data = ""
			}
			HH := 資料內容[75:77]
			MM := 資料內容[77:79]
			SS := 資料內容[79:81]
			sss := 資料內容[81:84]
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
			temp.data = "55688|" + temp.ttime + "|" + 資料內容
		} else if len(資料內容) == 162 { //emst委託
			if 資料內容[11:12] == "V" {
				temp.data = ""
			} else if 資料內容[84:85] == "E" {
				temp.data = "風控Error:" + 資料內容
				fmt.Println(temp.data)
			} else {
				HH := 資料內容[71:73]
				MM := 資料內容[73:75]
				SS := 資料內容[75:77]
				sss := 資料內容[77:80]
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
				temp.data = "55688|" + temp.ttime + "|" + 資料內容
			}
		} else if strings.Contains(資料內容, "order") { //Ben
			qqq := strings.Split(資料內容, ",")
			temp.ttime = qqq[0]
			if qqq[2] == "1t" {
				temp.data = "3333 |" + temp.ttime + "|" + "\x01" + qqq[3]
			} else if qqq[2] == "1s" {
				temp.data = "13334|" + temp.ttime + "|" + "\x01" + qqq[3]
			} else {
				temp.data = "13335|" + temp.ttime + "|" + "\x01" + qqq[3]
			}
		} else if strings.Contains(資料內容, "FIX.4.4") || strings.Contains(資料內容, "TMP_") {
			ooo := strings.Split(資料內容, ",")
			temp.ttime = ooo[0]
			if ooo[1] == "FIX_O" || ooo[1] == "FIX_T" {
				長度 := len(ooo[2])
				www := strings.Split(ooo[2], "\x01")
				for i := 0; i < len(www); i++ {
					if strings.Contains(www[i], "34=") {
						temp.流口水號 = www[i][3:]
					}
				}
				HH := temp.ttime[0:2]
				MM := temp.ttime[3:5]
				SS := temp.ttime[6:8]
				sss := temp.ttime[9:]
				// i, _ := strconv.Atoi(sss)
				// i = i + 40
				// if i >= 999 {
				// 	sss = strconv.Itoa(i - 999)
				// 	if len(sss) == 2 {
				// 		sss = "0" + sss
				// 	} else if len(sss) == 1 {
				// 		sss = "00" + sss
				// 	}
				// 	l, _ := strconv.Atoi(SS)
				// 	l = l + 1
				// 	if l >= 60 {
				// 		mmm, _ := strconv.Atoi(MM)
				// 		mmm = mmm + 1
				// 		MM = strconv.Itoa(mmm)
				// 		if mmm >= 60 {
				// 			hh, _ := strconv.Atoi(HH)
				// 			hh = hh + 1
				// 			HH = strconv.Itoa(hh)
				// 			if len(HH) < 2 {
				// 				HH = "0" + HH
				// 			}
				// 			MM = strconv.Itoa(mmm - 60)
				// 		}
				// 		if len(MM) < 2 {
				// 			MM = "0" + MM
				// 		}
				// 		SS = strconv.Itoa(l - 60)
				// 		if len(SS) < 2 {
				// 			SS = "0" + SS
				// 		}
				// 	} else {
				// 		SS = strconv.Itoa(l)
				// 		if len(SS) < 2 {
				// 			SS = "0" + SS
				// 		}
				// 	}
				// } else {
				// 	sss = strconv.Itoa(i)
				// 	if len(sss) == 2 {
				// 		sss = "0" + sss
				// 	} else if len(sss) == 1 {
				// 		sss = "00" + sss
				// 	}
				// }
				temp.ttime = HH + ":" + MM + ":" + SS + "." + sss
				temp.data = "55689|" + temp.ttime + "|" + ooo[2][:長度-8]
				temp.data = strings.ReplaceAll(temp.data, "body為:", "")
			} else if ooo[1] == "TMP_O" || ooo[1] == "TMP_T" {
				HH := temp.ttime[0:2]
				MM := temp.ttime[3:5]
				SS := temp.ttime[6:8]
				sss := temp.ttime[9:]
				// i, _ := strconv.Atoi(sss)
				// i = i + 45
				// if i >= 999 {
				// 	sss = strconv.Itoa(i - 999)
				// 	if len(sss) == 2 {
				// 		sss = "0" + sss
				// 	} else if len(sss) == 1 {
				// 		sss = "00" + sss
				// 	}
				// 	l, _ := strconv.Atoi(SS)
				// 	l = l + 1
				// 	if l >= 60 {
				// 		mmm, _ := strconv.Atoi(MM)
				// 		mmm = mmm + 1
				// 		MM = strconv.Itoa(mmm)
				// 		if mmm >= 60 {
				// 			hh, _ := strconv.Atoi(HH)
				// 			hh = hh + 1
				// 			HH = strconv.Itoa(hh)
				// 			if len(HH) < 2 {
				// 				HH = "0" + HH
				// 			}
				// 			MM = strconv.Itoa(mmm - 60)
				// 		}
				// 		if len(MM) < 2 {
				// 			MM = "0" + MM
				// 		}
				// 		SS = strconv.Itoa(l - 60)
				// 		if len(SS) < 2 {
				// 			SS = "0" + SS
				// 		}
				// 	} else {
				// 		SS = strconv.Itoa(l)
				// 		if len(SS) < 2 {
				// 			SS = "0" + SS
				// 		}
				// 	}
				// } else {
				// 	sss = strconv.Itoa(i)
				// 	if len(sss) == 2 {
				// 		sss = "0" + sss
				// 	} else if len(sss) == 1 {
				// 		sss = "00" + sss
				// 	}
				// }
				temp.ttime = HH + ":" + MM + ":" + SS + "." + sss
				temp.data = "55690|" + temp.ttime + "|" + ooo[2]
			}
		} else if len(資料內容) > 100 && 資料內容[84:85] == "E" {
			temp.data = "風控Error|" + 資料內容
			fmt.Println(temp.data)
		} else if strings.Contains(資料內容, "WAMT") {
			ooo := strings.Split(資料內容, ",")
			temp.ttime = ooo[0][9:] + ".000"
			a := strings.Index(資料內容, "'")
			b := strings.LastIndex(資料內容, "'")
			ddd := 資料內容[a+1 : b]
			ddd = strings.Replace(ddd, "\\", "", -1)
			temp.data = "3306 |" + temp.ttime + "|" + ddd
		} else if strings.Contains(資料內容, ".TXT") { //盤中異動
			temp.ttime = 資料內容[0:5]
			temp.ttime = temp.ttime + ":00.000"
			//指令 := "mv /home/Projects/file/" + 資料內容[6:] + "  /home/ftpinstant/"
			temp.data = "0000 |" + temp.ttime + "|" + 資料內容[6:]
		} else { //price
			m_delimiter := "\x01"
			qqq := strings.Split(資料內容, ",")
			temp.ttime = 資料內容[11:19]
			temp.ttime = temp.ttime + ".000"
			stock_id := strings.Replace(qqq[1], "\"", "", -1)
			電文 := "35=207" + m_delimiter + "55=" + stock_id + m_delimiter + "44=" + qqq[2] + m_delimiter
			temp.data = "8986 |" + temp.ttime + "|" + 電文
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
