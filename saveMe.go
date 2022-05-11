package main

// go doc  github.com/docker/docker/api/types  ImageSummary






import (
"archive/tar"
"compress/gzip"
"context"
"encoding/base64"
"encoding/json"
"flag"
"io"
"path/filepath"
//"errors"
//"path"
"time"

"github.com/tidwall/gjson"
"io/ioutil"
"os"
"strconv"

"fmt"

"github.com/docker/docker/api/types"
"github.com/docker/docker/client"
"strings"
)


func pullImageFromRepo(imageName ,host ,user,passwd string) {

	cli, err := client.NewClientWithOpts(client.WithHost("tcp://"+host), client.WithAPIVersionNegotiation())

	if err != nil {
		fmt.Println(err)
	}

	ctx := context.Background()

	authConfig := types.AuthConfig{
		Username: user,
		Password: passwd,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		panic(err)
	}

	defer out.Close()
	io.Copy(os.Stdout, out)

}

func main() {

	var image string
	var host string

	var username  string

	var password  string


	flag.StringVar(&image,"i","","镜像")

	flag.StringVar(&host,"h","","docker engine ip:port")

	flag.StringVar(&username,"u","","镜像仓库用户名")

	flag.StringVar(&password,"p","","镜像仓库密码")




	flag.Parse()

	fmt.Printf("image:%v\n",image)


	fmt.Printf("host:%v\n",host)


	imageName := image
	pullImageFromRepo(imageName,host,username,password)
	id := getImageId(imageName,host)
	if id == "" {
		fmt.Println("镜像不存在")
		return
	}

	strArry := strings.Split(imageName, "/")
	path := strArry[(len(strArry) - 1)]
	pathArray := strings.Split(path, ":")
	saveDir := pathArray[0] + "-" + pathArray[1]
	fmt.Println(saveDir)

	//
	done := make(chan struct{})
	defer close(done)

	f := func(done <-chan struct{}) {
		for {

			select {

			//case <-done:
			//	return
			default:
				time.Sleep(time.Second * 2)
				fmt.Printf(".")
				//fmt.Printf("Sec")

				//fmt.Fprintf(os.Stdout, "%s", ".")
				//fmt.Fprintf(os.Stdout, ".")
				//time.Sleep(time.Second * 1)
			case <-done:
				return

			}

		}

	}
	go f(done)

	//
	saveImage(id, saveDir,host)
	// 打开tar包  打开成功后删除原tar 包
	fmt.Println("从远程docker 引擎工作区拉取镜像元数据成功")
	fmt.Println("开始配置镜像元信息")
	untarFromPath(saveDir, "./"+saveDir+"/"+saveDir+".tar")
	err := os.Remove("./" + saveDir + "/" + saveDir + ".tar")
	if err != nil {
		fmt.Println(err)
	}
	// 配置元信息
	arry := strings.Split(imageName, ":")
	configMetdata(arry[0], arry[1], saveDir, imageName)
	// 开始打包
	//tarFun(saveDir+".tar", saveDir, saveDir)
	fmt.Println("镜像元信息配置成功开始打包镜像")
	tarFunWindows(saveDir+".tar", saveDir, saveDir)
	err = os.RemoveAll(saveDir)
	if err != nil {
		fmt.Println(err)
	}
}



func tarFunWindows(desc string, src string, dirSrc string) error {
	fd, err := os.Create(desc)
	if err != nil {
		return err
	}
	defer fd.Close()
	gw := gzip.NewWriter(fd)
	defer gw.Close()
	tr := tar.NewWriter(gw)
	defer tr.Close()
	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		strPath := strings.Replace(path, "\\", "/", -1)
		arrary := strings.Split(strPath, "/")
		if arrary[0] == dirSrc && path != dirSrc {
			hdr.Name = strings.Replace(strPath, dirSrc+"/", "", 1)
			err = tr.WriteHeader(hdr)
			if err != nil {
				return err
			}
			fs, err := os.Open(path)
			if err != nil {
				return err
			}
			defer fs.Close()
			if info.Mode().IsRegular() {
				io.Copy(tr, fs)
			}
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func configMetdata(imageName string, tag string, path string, RepoTags string) {
	data, err := ioutil.ReadFile(path + "/" + "manifest.json")
	if err != nil {
		return
	}
	str := string(data)
	length := gjson.Get(str, "0.Layers.#")
	n, _ := strconv.Atoi(length.String())
	n = n - 1
	//fmt.Println(n)
	//valur := gjson.Get(str, "0.Layers.6")

	strNum := strconv.Itoa(n)
	t := "0.Layers." + strNum
	valur1 := gjson.Get(str, t)

	//fmt.Println(length.String())
	//fmt.Println(valur.String())
	//fmt.Println(valur1.String())

	test := valur1.String()
	arr := strings.Split(test, "/")
	//fmt.Println(arr[0])
	lastId := arr[0]

	//
	//imageName := "nginx"
	//  {"nginx":{"latest":"9c77a26fdf5fb977da1d0e46e6c3cf5bb198f99d8c9b7d96118493b8fafe1d79"}}

	//tag := "latest"
	//
	str1 := "{" + `"` + imageName + `"` + `:{"` + tag + `":"` + lastId + `"}}`
	//	fmt.Println(str1)

	//fileName := "image/test.dat"
	fileName := path + "/" + "repositories"
	dstFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dstFile.Close()
	dstFile.WriteString(str1 + "\n")
	//
	//fmt.Println(str)
	//fmt.Println(strings.Replace(str, "null", "[\""+RepoTags+"\"]", 1))
	err = os.Remove(path + "/" + "manifest.json")
	if err != nil {
		panic(err)
	}

	//
	//imageName

	//manifestStr := strings.Replace(str, "null", "[\"172.16.100.216/binhu/lz-zyc-synergy-service:latest\"]", 1)
	manifestStr := strings.Replace(str, "null", "[\""+RepoTags+"\"]", 1)

	manifest, err := os.Create(path + "/" + "manifest.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer manifest.Close()
	manifest.WriteString(manifestStr + "\n")

}
func getImageId(tag string,host string) string {

	cli, err := client.NewClientWithOpts(client.WithHost("tcp://"+host), client.WithAPIVersionNegotiation())

	if err != nil {
		fmt.Println(err)
	}
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	for _, image := range images {
		tags := image.RepoTags
		for _, name := range tags {
			if name == tag {
				imageId := strings.Split(image.ID, ":")
				id := imageId[1]
				return id
			}
		}
	}
	return ""

}

func saveImage(id string, iname string,host string) {
	fmt.Println("开始从远程docker引擎工作区获取镜像")
	err := os.Mkdir(iname, 0666)
	if err == os.ErrExist {
		if err := os.Remove(iname); err != nil {
			fmt.Println("清除环境失败")
		} else {
			panic(err)
		}
	}
	fw, err := os.Create(iname + "/" + iname + ".tar")
	if err != nil {
		panic(err)
	}
	defer fw.Close()


	cli, err := client.NewClientWithOpts(client.WithHost("tcp://"+host), client.WithAPIVersionNegotiation())

	if err != nil {
		fmt.Println(err)
	}
	a := []string{id}
	r, err := cli.ImageSave(context.Background(), a)
	if err != nil {
		fmt.Println(err)
	}
	if _, err := io.Copy(fw, r); err != nil {
		panic(err)
	} else {
		//		fmt.Println(n)
	}
}

func makeDir() string {
	err := os.Mkdir("image", 0666)
	if err != nil {
		fmt.Println(err)
	}
	return "image"

}


func untarFromPath(base string, imageTar string) {
	r, err := os.Open(imageTar)
	if err != nil {
		panic(err)
	}
	defer r.Close()
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return
		}
		if err != nil {
			panic(err)
		}
		fullpath := filepath.Join(base, hdr.Name)
		info := hdr.FileInfo()
		// as dir
		if info.IsDir() {
			os.MkdirAll(fullpath, 0755)
			continue
		}
		dir := filepath.Dir(fullpath)
		os.MkdirAll(dir, 0755)
		// as file
		f, err := os.Create(fullpath)
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(f, tr)
		if err != nil {
			f.Close()
			panic(err)
		}
		f.Chmod(info.Mode())
		f.Close()
	}
}
