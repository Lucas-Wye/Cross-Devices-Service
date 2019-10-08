package models

import (
	"fmt"
	"io/ioutil"
)

var filelist = "./upload/filelist"

func GetFileList() []string {
	var res []string
	// fi, err := os.Open(filelist)
	// if err != nil {
	// 	fmt.Printf("Error: %s\n", err)
	// 	return nil
	// }
	// defer fi.Close()

	// br := bufio.NewReader(fi)
	// for {
	// 	a, _, c := br.ReadLine()
	// 	if c == io.EOF {
	// 		break
	// 	}
	// 	res = append(res, string(a))
	// }
	// return res

	pathname := "./upload"
	rd, err := ioutil.ReadDir(pathname)
	checkError(err)
	for _, fi := range rd {
		// if fi.IsDir() {
		// 	fmt.Printf("[%s]\n", pathname+"\\"+fi.Name())
		// 	GetAllFile(pathname + fi.Name() + "\\")
		// } else {
		// 	fmt.Println(fi.Name())
		// }
		res = append(res, string(fi.Name()))
	}
	fmt.Println(res)
	return res

}

// 上传的文件写入到清单中
// func WriteToFileList(filename string) bool {
// 	var wireteString = filename + "\n"
// 	var f *os.File
// 	var err1 error

// 	if checkFileIsExist(filelist) { //如果文件存在
// 		f, err1 = os.OpenFile(filelist, os.O_APPEND|os.O_WRONLY, os.ModeAppend) //打开文件
// 	} else {
// 		f, err1 = os.Create(filelist) //创建文件
// 		fmt.Println("文件", filelist, "不存在！")
// 	}
// 	checkError(err1)

// 	_, err1 = io.WriteString(f, wireteString) //写入文件(字符串)
// 	defer f.Close()

// 	// go tf_pose(filename)

// 	checkError(err1)
// 	if err1 != nil {
// 		return false
// 	}
// 	return true
// }

// 判断文件是否存在，存在返回true不存在返回false
// func checkFileIsExist(filename string) bool {
// 	var exist = true
// 	if _, err := os.Stat(filename); os.IsNotExist(err) {
// 		exist = false
// 	}
// 	return exist
// }

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// func tf_pose(infile string) bool {
// 	// tf_pose处理
// 	outfile := "out_" + infile
// 	cmd := exec.Command("python", "run.py --model=mobilenet_thin  --resize=432x368", infile, outfile)
// 	out, err := cmd.CombinedOutput()
// 	checkError(err)
// 	fmt.Println(string(out))
// 	if err != nil {
// 		return false
// 	}
// 	return true
// }
