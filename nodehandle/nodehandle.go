package nodehandle

import (

	// go
	//"log"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	// local
	"gnvm/config"
)

const (
	DIVIDE = "\\"
	NODE   = "node.exe"
)

var globalNodePath, rootPath string

func GetGlobalNodePath() string {

	// get node.exe path
	file, err := exec.LookPath(NODE)
	if err != nil {
		globalNodePath = "root"
	} else {
		// relpace "\\node.exe"
		globalNodePath = strings.Replace(file, DIVIDE+NODE, "", -1)
	}

	// gnvm.exe and node.exe the same path
	if globalNodePath == "." {
		globalNodePath = "root"
	}
	//log.Println("Node.exe path is: ", globalNodePath)

	return globalNodePath
}

func getCurrentPath() string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Get current path Error: " + err.Error())
		return ""
	}
	return path
}

func isDirExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		// return file.IsDir()
		return true
	}
}

func getNodeVersion(path string) (string, error) {
	var newout string
	out, err := exec.Command(path+"node", "--version").Output()
	//string(out[:]) bytes to string
	if err == nil {
		// replace \r\n
		newout = strings.Replace(string(string(out[:])[1:]), "\r\n", "", -1)
	}
	return newout, err
}

func cmd(name, arg string) error {
	_, err := exec.Command("cmd", "/C", name, arg).Output()
	return err
}

func copy(src, dest string) error {
	_, err := exec.Command("cmd", "/C", "copy", "/y", src, dest).Output()
	return err
}

/**
 * rootPath is gnvm.exe root path,     e.g <root>
 */
func SetRootPath() {

	// set rootPath and rootNode
	if globalNodePath == "root" {
		rootPath = getCurrentPath() + DIVIDE
	} else {
		rootPath = globalNodePath + DIVIDE
	}
	//log.Println("Current path is: " + rootPath)
}

/**
 * rootNode is rootPath + "/node.exe", e.g. <root>/node.exe
 *
 * usePath  is use node version path,  e.g. <root>/x.xx.xx
 * useNode  is usePath + "/node.exe",  e.g. <root>/x.xx.xx/node.exe
 *
 * rootVersion is <root>/node.exe version
 * rootFolder  is <root>/rootVersion
 */
func Use(folder string) bool {

	// get true folder, e.g. folder is latest return x.xx.xx
	folder = GetTrueVersion(folder)

	if folder == config.UNKNOWN {
		fmt.Println("Waring: Unassigned Node.js latest version. See 'gnvm install latest'.")
		return false
	}

	// set rootNode
	rootNode := rootPath + NODE
	//log.Println("Root node path is: " + rootNode)

	// set usePath and useNode
	usePath := rootPath + folder + DIVIDE
	useNode := usePath + NODE
	//log.Println("Use node.exe path is: " + usePath)

	// get <root>/node.exe version
	rootVersion, err := getNodeVersion(rootPath)
	if err != nil {
		fmt.Println("Not found global node version, please checkout. If not exist node.exe. See 'gnvm install latest'.")
		return false
	}
	//log.Printf("Root node.exe verison is: %v \n", rootVersion)

	// <root>/folder is exist
	if isDirExist(usePath) != true {
		fmt.Printf("[%v] folder is not exist. Get local node.exe version. See 'gnvm ls'.", folder)
		return false
	}

	// <root>/node.exe is exist
	if isDirExist(rootNode) != true {
		fmt.Println("Not found global node version, please checkout. If not exist node.exe. See 'gnvm install latest'.")
		return false
	}

	// check folder is rootVersion
	if folder == rootVersion {
		fmt.Printf("Current node.exe version is [%v], not re-use. See 'gnvm node-version'.", folder)
		return false
	}

	// set rootFolder
	rootFolder := rootPath + rootVersion

	// <root>/rootVersion is exist
	if isDirExist(rootFolder) != true {

		// create rootVersion folder
		if err := cmd("md", rootFolder); err != nil {
			fmt.Printf("Create %v folder Error: %v", rootVersion, err.Error())
			return false
		}

	}

	// copy rootNode to <root>/rootVersion
	if err := copy(rootNode, rootFolder); err != nil {
		fmt.Printf("copy %v to %v folder Error: ", rootNode, rootFolder, err.Error())
		return false
	}

	// delete <root>/node.exe
	if err := os.Remove(rootNode); err != nil {
		fmt.Printf("remove %v to %v folder Error: ", rootNode, err.Error())
		return false
	}

	// copy useNode to rootPath
	if err := copy(useNode, rootPath); err != nil {
		fmt.Printf("copy %v to %v folder Error: ", useNode, rootPath, err.Error())
		return false
	}

	fmt.Printf("Set success, Current Node.exe version is [%v].", folder)

	return true

}

func VerifyNodeVersion(version string) bool {
	result := true
	if version == config.UNKNOWN {
		return true
	}
	arr := strings.Split(version, ".")
	if len(arr) != 3 {
		return false
	}
	for _, v := range arr {
		_, err := strconv.ParseInt(v, 10, 0)
		if err != nil {
			result = false
			break
		}
	}
	return result
}

func GetTrueVersion(latest string) string {
	if latest == config.LATEST {
		latest = config.GetConfig(config.LATEST_VERSION)
		fmt.Printf("Current latest version is [%v] \n", latest)
	}
	return latest
}

func NodeVersion() {
	latest := config.GetConfig(config.LATEST_VERSION)
	global := config.GetConfig(config.GLOBAL_VERSION)

	fmt.Printf(`Node.exe global verson is [%v]
Node.exe latest verson is [%v]
Notice: When version is [%v], please See 'gnvm use help'.`, global, latest, config.UNKNOWN)
}

func Uninstall(folder string) {

	// set removePath
	removePath := rootPath + folder

	if folder == config.UNKNOWN {
		fmt.Println("Waring: Unassigned Node.js latest version. See 'gnvm install latest'.")
		return
	}

	// rootPath/version is exist
	if isDirExist(removePath) != true {
		fmt.Printf("Waring: [%v] folder is not exist. Get local node.exe version. See 'gnvm ls'.\n", folder)
		return
	}

	// remove rootPath/version folder
	if err := os.RemoveAll(removePath); err != nil {
		fmt.Printf("Uinstall [%v] fail, Error: %v", folder, err.Error())
		return
	}

	fmt.Printf("Node.exe version [%v] uninstall success. \n", folder)
}