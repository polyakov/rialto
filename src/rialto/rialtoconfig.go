package rialto

import "fmt"
import "encoding/json"
//import "os"
import "bytes"
//import "bufio"
import "io/ioutil"
//import "net/http"
import "text/template"
//import "github.com/Masterminds/sprig"
import "strings"

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func LoadValueTemplate(serviceInstance ServiceInstance) string {
	var url string = serviceInstance.ChartSource.Url+
		"/"+serviceInstance.ChartSource.ChartName+
		"/values.template.yaml"
	fmt.Println(url)

	return strings.TrimSpace( url )

}

// eventully we want to be able to say "I need a FHIR server" and have the system find as service to fit this requirement
/*
func ReadWellKnownServiceConfig( filepath string ) WellKnownServiceConfig {
    fmt.Printf("filename: %v", filepath)

    dat, err := ioutil.ReadFile(filepath)
    check(err)

    //var jsonStr = string(dat)
    //fmt.Print(jsonStr)

    res := WellKnownServiceConfig{}
    json.Unmarshal(dat, &res)
    fmt.Println(res)

    return res
}
*/

func ApplyTemplate( filepath string, deployment Deployment)[]byte {
	dat, err := ioutil.ReadFile(filepath)
	check(err)
	var valStr string = string(dat)

	tmpl, err := template.New("value").Parse(valStr)
	if err != nil { panic(err) }
	var outBuf bytes.Buffer;
	//err = tmpl.Execute(os.Stdout, deployment)
	err = tmpl.Execute(&outBuf, deployment)
	if err != nil { panic(err) }
	//fmt.Println("==============")
	//fmt.Println(outBuf)
	return outBuf.Bytes()
}


func GetServiceInstance(foundServiceConfig ServiceConfig, name string) ServiceInstance {
    var serviceInstance ServiceInstance=ServiceInstance{}
    serviceInstance.Name=name
    serviceInstance.ServiceName=foundServiceConfig.Name
    serviceInstance.Exposes=map[string]ServiceInstanceProperties{}
    // map Exposes
    for interfaceIndex,foundInterface:= range foundServiceConfig.Exposes {
	    fmt.Printf("Found interface: %v: %v\n", interfaceIndex, foundInterface.Name)
	    var exposedProperties ServiceInstanceProperties = ServiceInstanceProperties{}
	    exposedProperties.Properties=map[string]string{}
	    for indexProp,foundProperty:= range foundInterface.Properties {
		    fmt.Printf("Found prop: %v: %v\n", indexProp, foundProperty.Name)
		    exposedProperties.Properties[foundProperty.Name]="??";
	    }
	    serviceInstance.Exposes[foundInterface.Name] = exposedProperties
    }

    return serviceInstance;
}

//ReadServiceConfig - parse JSON service config and return ServiceConfig stuct
func ReadServiceConfig( filepath string)  ServiceConfig {
    //fmt.Printf("filename: %v", filepath)

    dat, err := ioutil.ReadFile(filepath)
    check(err)

    //var jsonStr = string(dat)
    //fmt.Print(jsonStr)

    res := ServiceConfig{}
    json.Unmarshal(dat, &res)
    //fmt.Println(res)

    return res;
}

//ReadInstallConfig - read a JSON doc with configuration of service instances
func ReadInstallConfig( filepath string)  InstallConfig {
    //fmt.Printf("filename: %v", filepath)

    dat, err := ioutil.ReadFile(filepath)
    check(err)

    //var jsonStr = string(dat)
    //fmt.Print(jsonStr)

    res := InstallConfig{}
    json.Unmarshal(dat, &res)
    fmt.Println(res)

    return res;
}
