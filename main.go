package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// TODO 需要忽略隐藏文件、系统文件等

func walkFunc(path string, info os.FileInfo, err error) error {
	if info == nil {
		fmt.Println("没发现目录或文件：", path) // 目录或文件名太长也会报错
		return nil
	}
	sqlStr := "insert into files(path,hash,dis) values (?,?,?)"

	_, err = Mydb.Exec(sqlStr, path, gethash(path), dis)
	if err != nil {
		fmt.Println("录入数据库失败：", err)
		return err
	}
	//theID, err := ret.LastInsertId() // 新插入数据的id
	//if err != nil {
	//	fmt.Println("取ID失败: ", err)
	//	return err
	//}
	//fmt.Println("保存成功：", theID)

	// time.Sleep(1 * time.Nanosecond)
	// if info.IsDir() {
	// fmt.Println("是个文件夹: ", path)
	//	return nil
	//} else {
	// fmt.Println("是个文件: ", path)
	//	return nil
	// }
	return nil
}

func findFileList(root string) {
	err := filepath.Walk(root, walkFunc)
	if err != nil {
		fmt.Println("filepath.Walk() 出错：", err)
	}
	return
}

func NewDB() *sql.DB {
	DB, _ := sql.Open("mysql", "root:root@tcp(localhost:3306)/myfile?charset=utf8")

	DB.SetConnMaxLifetime(150 * time.Second)
	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(10)

	if err := DB.Ping(); err != nil {
		fmt.Println("连接数据库失败")
		return nil
	}
	fmt.Println("数据库连接成功")
	return DB
}

func gethash(path string) (hash string) {
	file, _ := os.Open(path)
	h_ob := sha256.New()
	_, err := io.Copy(h_ob, file)
	if err == nil {
		hash := h_ob.Sum(nil)
		hashvalue := hex.EncodeToString(hash)
		return hashvalue
	} else {
		return ""
	}
}

func start(no int, c chan string, e chan bool) {
	fmt.Println("GO程", no, "启动")
	timeStart := time.Now()
	for {
		data, ok := <-c

		if !ok {
			timeEnd := time.Since(timeStart)
			fmt.Println("GO程", no, "退出，耗时", timeEnd)
			e <- true
			break
		}

		fmt.Println("GO程", no, "开始遍历文件夹", data)
		findFileList(data)
	}
}

var Mydb = NewDB()
var dis string

func main() {

	timeStart := time.Now()

	count := 10 // 启动多少个go程
	c := make(chan string, count)
	e := make(chan bool, count) //用于标识go程退出

	var path string
	dis = "music"

	fmt.Println("请输入想要遍历的路径：")
	fmt.Scanf("%s", &path)
	fmt.Println("=====Start=====")

	for i := 0; i < count; i++ {
		go start(i, c, e)
	}

	files, _ := ioutil.ReadDir(path)
	// TODO 用map保存未遍历的目录
	// m := make(map[int]string)

	for _, info := range files {
		if info.IsDir() {
			// go showFileList(filepath.Join(path, info.Name()))
			c <- filepath.Join(path, info.Name())
		} else {
			// fmt.Printf(filepath.Join(path, info.Name()) + "\n")
			// TODO 写数据库
		}
	}

	close(c)

	for i := 0; i < count; i++ {
		<-e
	}
	//for {
	//	if len(e) == count {
	//		break
	//	}
	//}
	timeEnd := time.Since(timeStart)
	fmt.Println(timeEnd)
	fmt.Println("主线程退出")
}
